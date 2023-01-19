package controllers

import (
	"fmt"
	_ "fmt"
	"mirauserlab/models"
	"mirauserlab/utils/tokens"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CurrentUser godoc
// @Summary  Gives the information about the current user
// @Schemes
// @Description  Gives the information about the current user.
// @Tags         write
// @Accept       json
// @Produce      json
// @Success      200
// @Security ApiKeyAuth
// @in header
// @name ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router       /api/admin/currentuser [get]

func CurrentUser(c *gin.Context) {

	user_id, err := tokens.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

// LoginInput is a struct with a subset of the fields of controllers.LoginInput. It is used when
// LoginInput needs to be provided as an input for data insertion from database. It does not include
// auto-generated fields.
type LoginInput struct {
	Username string `json:"username" binding:"required,endswith=@mirafra.com"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary  Logs in the user and generates token
// @Schemes
// @Description  Logs in the user and generates token.
// @Tags         write
// @Accept       json
// @Produce      json
// @Param        input body  LoginInput  true  "New User"
// @Success      200
// @Router       /api/login [post]
func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Bearer Token": token})

}

// Login godoc
// @Summary  Get the current user
// @Schemes
// @Description  Get the information about the current user.
// @Tags         write
// @Accept       json
// @Produce      json
// @Success      200
// @Security BearerAuth
// @in header
// @name Authorization
// @param Authorization header string true "Authorization"
// @TOKEN
// @Router       /api/userget [get]
func UserGet(c *gin.Context) {
	//Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTk2MTM5MTgsInVzZXJfaWQiOjh9.bItdHSo-YZjAEZfa1-wuunpD1Y4WTEilsFnbcPT13t8
	user_id, err := tokens.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

// RegisterInput is a struct with a subset of the fields of controllers.RegisterInput. It is used when
// RegisterInput needs to be provided as an input for data insertion from json body. It does not include
// auto-generated fields.
type RegisterInput struct {
	Username string `json:"username" binding:"required,endswith=@mirafra.com"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary  Registers new user
// @Schemes
// @Description  Registers new user.
// @Tags         write
// @Accept       json
// @Produce      json
// @Param        input body  RegisterInput  true  "New User"
// @Success      200
// @Router       /api/register [post]
func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid Username or password submitted!": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while saving user in database!": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username and password successfully registered in database!"})
}

type InputUsers struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Title : Create Password
// @Description: Creates a new password for the specified user
// @Header: main header
// @Param request body main.CreatePassword true "Generate a new password"
// security:
// - apiKey: []
// responses:
// @Success 200: main.CreatePasswordResponse
// @Router /createpassword [patch]
func CreatePassword(c *gin.Context) {

	var user models.User
	if err := models.DB.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Record Not found"})
		return
	}

	var input InputUsers
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	models.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data updated successfully!", "data": user})
}

// @Title : Get all users
// @Description: Gets all users present inside the database
// @Header: main header
// @Param request body main.GetAllUsers true "Gives all users in the database"
// security:
// - apiKey: []
// responses:
// @Success 200: main.GetAllUsersResponse
// @Router /getallusers [get]
func GetAllUsers(c *gin.Context) {

	var user []models.User
	if err := models.DB.Find(&user).Error; err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data retrieved successfully", "data": user})
}

type CreateEmployeeInput struct {
	EmpId   int    `json:"empId"`
	EmpName string `json:"empName"`
	Address string `json:"address"`
	Phone   int64  `json:"phone"`
	AssetId int64  `json:"assetId"`
}

func CreateEmployee(c *gin.Context) {

	var input CreateEmployeeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	if input.EmpName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"stauts": false, "message": "Please fill all the fields!"})
		return

	}
	if input.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"stauts": false, "message": "Please fill all the fields!"})
		return
	}
	emp := models.Employee{
		EmpId:   input.EmpId,
		EmpName: input.EmpName,
		Address: input.Address,
		Phone:   input.Phone,
		AssetId: input.AssetId,
	}
	if err := models.DB.Create(&emp).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Duplicate entry found!"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": true, "message": "Employee data inserted successfully!", "data": emp})

}

func GetAssetById(c *gin.Context) {

	var asset models.Asset

	if err := models.DB.Where("asset_id=?", c.Param("id")).First(&asset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data retrieved successfully!", "data": asset})

}

type InputAssets struct {
	AssetId          int64  `json:"assetId"`
	AssetDescription string `json:"assetDescription"`
}

func CreateAssets(c *gin.Context) {
	var input InputAssets

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
	}

	if input.AssetDescription == "" {
		c.JSON(http.StatusBadRequest, gin.H{"stauts": false, "message": "Please fill all the fields!"})
		return
	}
	assets := models.Asset{
		AssetId:          input.AssetId,
		AssetDescription: input.AssetDescription,
	}

	if err := models.DB.Create(&assets).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Duplicate entry found!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": true, "message": "Assets data inserted successfully!", "data": assets})

}

func DeleteAssetById(c *gin.Context) {

	var asset models.Asset
	if err := models.DB.Where("asset_id=?", c.Param("id")).First(&asset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "No Record found!"})
		return
	}
	models.DB.Delete(&asset)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Asset data deleted successfully!"})
}

func GetAllAssets(c *gin.Context) {
	var asset []models.Asset
	models.DB.Find(&asset)

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data retrieved successfully", "data": asset})
}

type UpdateAssetIdInput struct {
	AssetId int64 `json:"assetId"`
}

func UpdateEmployeeForAssetId(c *gin.Context) {

	var emp models.Employee

	if err := models.DB.Where("emp_id=?", c.Param("id")).First(&emp).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "record not found!"})
		return
	}

	models.Connection().Exec("update Employees set asset_id=? where emp_id=?", c.Param("assetId"), c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data updatated successfully!"})

}

func DeRegisterByAssetId(c *gin.Context) {
	var emp models.Employee

	if err := models.DB.Where("emp_id=?", c.Param("id")).First(&emp).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "record not found!"})
		return
	}

	models.Connection().Exec("update Employees set asset_id=? where emp_id=?", "", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "data updatated successfully!"})

}

func Logout(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
		c.Abort()
		return
	}
	extractedToken := strings.Split(token, "Bearer ")
	fmt.Println(extractedToken)

	// claims := jwt.MapClaims{}
	// claims["authorized"] = true
	// claims["exp"] = time.Now().Add(time.Second * time.Duration(0)).Unix()
	// jwt.NewWithClaims(extractedToken, claims)

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Sucessfully logged out"})
}
