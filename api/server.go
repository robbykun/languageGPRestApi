package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/robbykun/languageGPRestApi/db"
)

func Init() {
	fmt.Println("api.Init開始")

	r := gin.Default()

	//DELETE table all data
	r.DELETE("/project", func(c *gin.Context) {
		db.GetDB().Delete(db.Project{})
		db.GetDB().Delete(db.Language{})
		c.JSON(http.StatusOK, "")
	})

	//CREATE
	r.POST("/project", func(c *gin.Context) {

		type ProjectJsonBody struct {
			Projects []struct {
				ProjectNo string `json:"project_no"`
				Price     int    `json:"price" binding:"min=0"`
				Station   string `json:"station"`
				Languages []struct {
					ProjectNo    string `json:"project_no"`
					LanguageType string `json:"language_type"`
				} `json:"languages"`
			} `json:"projects"`
		}

		var one_page_projects ProjectJsonBody

		if err := c.ShouldBindJSON(&one_page_projects); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, only_project := range one_page_projects.Projects {

			project := &db.Project{}
			copyStruct(only_project, project)

			fmt.Print("project")
			fmt.Printf("%+v\n", *project)
			db.GetDB().Create(&project)

			for _, only_language := range only_project.Languages {
				language := &db.Language{}

				copyStruct(only_language, language)

				fmt.Printf("%+v\n", language)
				db.GetDB().Create(&language)
			}
		}
		c.JSON(http.StatusOK, "")
	})

	// プログラミング言語ごとに集計
	r.GET("/group/language_type", func(c *gin.Context) {

		type Results struct {
			LanguageType string
			Avg          float64
			Max          uint
			Min          uint
			Count        uint
		}

		results := &[]Results{}

		db.GetDB().Table("projects").
			Select("language_type, avg(price) as \"avg\", max(price) as \"max\", min(price) as \"min\", count(*) as \"count\"").
			Joins("inner join languages on projects.project_no = languages.project_no").
			Group("language_type").
			Order("avg desc").
			Scan(&results)

		c.JSON(http.StatusOK, *results)
	})

	// 最寄駅ごとに集計
	r.GET("/group/station", func(c *gin.Context) {

		type Results struct {
			Station string
			Avg     float64
			Max     uint
			Min     uint
			Count   uint
		}

		results := &[]Results{}

		db.GetDB().Table("projects").
			Select("station, avg(price) as \"avg\", max(price) as \"max\", min(price) as \"min\", count(*) as \"count\"").
			Joins("inner join languages on projects.project_no = languages.project_no").
			Group("station").
			Order("count desc").
			Scan(&results)

		c.JSON(http.StatusOK, *results)
	})

	r.Run(":8080")
}

func copyStruct(src interface{}, dst interface{}) error {
	fv := reflect.ValueOf(src)

	ft := fv.Type()
	if fv.Kind() == reflect.Ptr {
		ft = ft.Elem()
		fv = fv.Elem()
	}

	tv := reflect.ValueOf(dst)
	if tv.Kind() != reflect.Ptr {
		return fmt.Errorf("[Error] non-pointer: %v", dst)
	}

	num := ft.NumField()
	for i := 0; i < num; i++ {
		field := ft.Field(i)

		if !field.Anonymous {
			name := field.Name
			srcField := fv.FieldByName(name)
			dstField := tv.Elem().FieldByName(name)

			if srcField.IsValid() && dstField.IsValid() {
				if srcField.Type() == dstField.Type() {
					dstField.Set(srcField)
				}
			}
		}
	}

	return nil
}