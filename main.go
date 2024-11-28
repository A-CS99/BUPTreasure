package main

import (
	"BUPTreasure/internal/ApiDTO"
	"BUPTreasure/internal/myDB"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type SignInfo = ApiDTO.SignInfo
type PickDTO = ApiDTO.PickDTO
type AvatarsDTO = ApiDTO.AvatarsDTO
type AssignDTO = ApiDTO.AssignDTO
type AllDTO = ApiDTO.AllDTO

var AwardTypes = ApiDTO.AwardTypes

func main() {
	// 设置指定获奖者
	var assigned []AssignDTO
	// 报名开始状态
	signStart := false

	// 初始化数据库
	err := myDB.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := myDB.Disconnect()
		if err != nil {
			fmt.Println(err)
		}
	}()
	r := gin.Default()
	// 设置跨域中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		// 如果是OPTIONS请求，直接返回200状态码，不再继续后续处理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// 继续处理请求
		c.Next()
	})
	// 移动端报名接口
	signApp := r.Group("/signApp")
	signApp.PUT("/flush", func(c *gin.Context) {
		// 清空User表 (删除User表并重新创建)
		err := myDB.FlushTable()
		if err != nil {
			fmt.Println(err)
			return
		} else {
			signStart = false
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": nil,
				"msg":  "Success Drop Table",
			})
		}
	})

	signApp.POST("/", func(c *gin.Context) {
		// 获取并存储报名信息
		signInfo := myDB.SignInfo{}
		err := c.ShouldBindJSON(&signInfo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "Bad Request",
			})
			return
		} else {
			if signStart {
				fmt.Println(signInfo.Name)
				fmt.Println(signInfo.AvatarUrl)
				err = myDB.Insert(signInfo)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"data": nil,
						"msg":  err,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": signInfo,
					"msg":  "OK",
				})
			} else {
				c.JSON(http.StatusAccepted, gin.H{
					"code": 202,
					"data": nil,
					"msg":  "Sign Not Start",
				})
			}
		}
	})

	signApp.GET("/on", func(c *gin.Context) {
		// 开始报名
		signStart = true
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": nil,
			"msg":  "Start Sign",
		})
	})

	signApp.GET("/off", func(c *gin.Context) {
		// 结束报名
		signStart = false
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": nil,
			"msg":  "Stop Sign",
		})
	})

	signApp.GET("/check", func(c *gin.Context) {
		// 查询是否已经报名
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "Bad Request: name is empty",
			})
			return
		}
		signInfo, err := myDB.QueryByName(name)
		if err != nil {
			fmt.Println("查询失败: ")
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"data": nil,
				"msg":  "No One Found",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": signInfo,
				"msg":  "OK",
			})
		}
	})

	// 网页端抽奖接口
	pickWeb := r.Group("/pickWeb")

	pickWeb.GET("/pick", func(c *gin.Context) {
		// 抽取中奖者
		var picked []SignInfo
		pickNumStr := c.DefaultQuery("pickNum", "0")
		pickNum, err := strconv.Atoi(pickNumStr)
		if err != nil {
			fmt.Println("数值字符串转换失败: ")
			fmt.Println(err)
			return
		}
		awardTypeIdxStr := c.Query("awardType")
		if awardTypeIdxStr == "" {
			fmt.Println("奖项类型不能为空")
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "Bad Request: wrong awardType",
			})
			return
		}
		awardTypeIdx, err := strconv.Atoi(awardTypeIdxStr)
		if err != nil {
			fmt.Println("数值字符串转换失败: ")
			fmt.Println(err)
			return
		}
		if awardTypeIdx < 0 || awardTypeIdx >= len(AwardTypes) {
			fmt.Println("奖项类型错误")
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "Bad Request: wrong awardType",
			})
			return
		}
		awardType := AwardTypes[awardTypeIdx]
		originPickNum := pickNum
		if len(assigned) > 0 {
			// 从指定中奖者中抽取
			assignedSlice := assigned
			if len(assigned) > pickNum {
				assignedSlice = assigned[:pickNum]
			}
			for _, assign := range assignedSlice {
				if assign.Award == awardType {
					prePick, err := myDB.QueryByName(assign.Name)
					if err != nil {
						fmt.Println("查询失败: ")
						fmt.Println(err)
					} else {
						picked = append(picked, prePick)
						pickNum--
					}
				}
			}
			assigned = []AssignDTO{}
			if pickNum == 0 {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": PickDTO{PickNum: originPickNum, AwardType: awardType, Picked: picked},
					"msg":  "OK",
				})
				return
			}
		}
		realPicked, err := myDB.RandomQuery(pickNum, awardType)
		if err != nil {
			fmt.Println("抽取失败: ")
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"data": nil,
				"msg":  "No One Found",
			})
			return
		}
		picked = append(picked, realPicked...)
		actualPickNum := len(picked)
		if actualPickNum == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"data": nil,
				"msg":  "No One Found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": PickDTO{PickNum: actualPickNum, AwardType: awardType, Picked: picked},
			"msg":  "OK",
		})
	})

	pickWeb.PUT("/flushPick", func(c *gin.Context) {
		// 清除中奖状况
		err := myDB.FlushAllAward()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"data": nil,
				"msg":  err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": nil,
			"msg":  "OK",
		})
	})

	pickWeb.GET("/total", func(c *gin.Context) {
		// 查询报名总人数
		total, err := myDB.TableSize()
		if err != nil {
			fmt.Println("查询失败: ")
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": total,
			"msg":  "OK",
		})
	})

	pickWeb.GET("/avatars", func(c *gin.Context) {
		// 查询特定范围内的头像
		fromStr := c.DefaultQuery("from", "0")
		toStr := c.DefaultQuery("to", "10")
		from, err := strconv.Atoi(fromStr)
		if err != nil {
			fmt.Println("数值字符串转换失败: ")
			fmt.Println(err)
			return
		}
		to, err := strconv.Atoi(toStr)
		if err != nil {
			fmt.Println("数值字符串转换失败: ")
			fmt.Println(err)
			return
		}
		avatarUrls, err := myDB.RangeQuery(from, to)
		if err != nil {
			fmt.Println("查询失败: ")
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"data": AvatarsDTO{AvatarNum: 0, AvatarUrls: nil},
				"msg":  "No One Found",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": AvatarsDTO{AvatarNum: len(avatarUrls), AvatarUrls: avatarUrls},
				"msg":  "OK",
			})
		}
	})

	pickWeb.GET("/all", func(c *gin.Context) {
		// 查询所有报名信息
		users, err := myDB.QueryAll()
		if err != nil {
			fmt.Println("查询失败: ")
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"data": nil,
				"msg":  "No One Found",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": AllDTO{UserNum: len(users), Users: users},
				"msg":  "OK",
			})
		}
	})

	pickWeb.PUT("/Ab2sRIjgFNo", func(c *gin.Context) {
		// 设置下一位中奖者
		name := c.Query("name")
		awardStr := c.Query("award")
		if name == "" || awardStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "Bad Request: name or award is empty",
			})
			return
		}
		award, err := strconv.Atoi(awardStr)
		if err != nil {
			fmt.Println("数值字符串转换失败: ")
			fmt.Println(err)
			return
		}
		AwardType := AwardTypes[award]
		assigned = append(assigned, AssignDTO{Name: name, Award: AwardType})
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": nil,
			"msg":  "OK",
		})
	})

	// 8080端口运行服务端程序
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Print(err)
		return
	}
}
