package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	// CREATE
	// 駅マスタ作成
	r.POST("/station", func(c *gin.Context) {

		// 路線情報を取得
		url := "http://express.heartrails.com/api/json?method=getLines&prefecture=%E6%9D%B1%E4%BA%AC%E9%83%BD"

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		type LineInfo struct {
			Response struct {
				Line []string `json:"line"`
			}
		}
		var lineInfo LineInfo

		json.Unmarshal(body, &lineInfo)

		// 駅情報を取得
		for _, line := range lineInfo.Response.Line {
			url = "http://express.heartrails.com/api/json?method=getStations&line=" + line + "&prefecture=%E6%9D%B1%E4%BA%AC%E9%83%BD"

			response, err = http.Get(url)
			if err != nil {
				log.Fatal(err)
			}

			body, err = ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			type StationInfo struct {
				Response struct {
					Station []struct {
						StationName string  `json:"name"`
						Keido       float64 `json:"x"`
						Ido         float64 `json:"y"`
					} `json:"station"`
				} `json:"response"`
			}
			var stationInfo StationInfo
			json.Unmarshal(body, &stationInfo)

			// DB登録
			for _, resStation := range stationInfo.Response.Station {

				station := &db.Station{}

				copyStruct(resStation, station)

				fmt.Println(station)

				db.GetDB().Create(&station)
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
			Ido     float64
			Keido   float64
			Avg     float64
			Max     uint
			Min     uint
			Count   uint
		}

		results := &[]Results{}

		db.GetDB().Debug().Table("projects").
			Select("station, ido, keido, avg(price) as \"avg\", max(price) as \"max\", min(price) as \"min\", count(*) as \"count\"").
			Joins("inner join languages on projects.project_no = languages.project_no").
			Joins("inner join stations on projects.station = stations.station_name").
			Group("station, ido, keido").
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
