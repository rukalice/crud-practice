package main

import "github.com/gin-gonic/gin"
import "net/http"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "log"

func main() {
    db, err := sql.Open("sqlite3", "./crud.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")

    r.StaticFile("/favicon.ico", "templates")

    r.GET("/", func(c *gin.Context) {
        users := make(map[int]string)
        rows, err := db.Query("select id, name from users")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()
        for rows.Next() {
            var id int
            var name string
            err = rows.Scan(&id, &name)
            if err != nil {
                log.Fatal(err)
            }
            users[id] = name
        }
        err = rows.Err()
        if err != nil {
            log.Fatal(err)
        }
        c.HTML(http.StatusOK, "index.tmpl", gin.H {
            "users": users,
        })
    })

    r.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
println(id)
        stmt, err := db.Prepare("select name from users where id = ?")
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()
        var name string
        err = stmt.QueryRow(id).Scan(&name)
        if err != nil {
            log.Fatal(err)
        }
        c.HTML(http.StatusOK, "detail.tmpl", gin.H {
            "name": name,
        })
    })
    
    r.POST("/", func(c *gin.Context) {
        name := c.PostForm("name")

        tx, err := db.Begin()
        if err != nil {
            log.Fatal(err)
        }

        var maxid int
        err = tx.QueryRow("select max(id) from users").Scan(&maxid)
        if err != nil {
            log.Fatal(err)
        }
        maxid++

        stmt, err := tx.Prepare("insert into users values(?, ?)")
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        _, err = stmt.Exec(maxid, name)
        if err != nil {
            log.Fatal(err)
        }
        err = tx.Commit()
        if err != nil {
            log. Fatal(err)
        }

        c.HTML(http.StatusOK, "detail.tmpl", gin.H {
            "id": maxid,
            "name": name,
        })
    })

    r.POST("/:id/update", func(c *gin.Context) {
        id := c.Param("id")
        name := c.PostForm("name")

        tx, err := db.Begin()
        if err != nil {
            log.Fatal(err)
        }

        stmt, err := tx.Prepare("update users set name = ? Where id = ?")
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        _, err = stmt.Exec(name, id)
        if err != nil {
            log.Fatal(err)
        }

        err = tx.Commit()
        if err != nil {
            log.Fatal(err)
        }

        c.HTML(http.StatusOK, "updated.tmpl", gin.H {
            "name": name,
        })
    })

    r.POST("/:id/delete", func(c *gin.Context) {
        id := c.Param("id")
println("delete")
        tx, err := db.Begin()
        if err != nil {
            log.Fatal(err)
        }

        stmt, err := tx.Prepare("delete from users where id = ?")
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        _, err = stmt.Exec(id)
        if err != nil {
            log.Fatal(err)
        }

        err = tx.Commit()
        if err != nil {
            log.Fatal(err)
        }

        c.HTML(http.StatusOK, "deleted.tmpl", gin.H {
            "id": id,
        })
    })

    r.Run()
}

