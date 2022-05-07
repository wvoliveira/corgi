package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elga-io/corgi/internal/app/config"
	"github.com/elga-io/corgi/internal/app/entity"
	"github.com/elga-io/corgi/internal/pkg/database"
	"github.com/google/uuid"
	"log"
)

var (
	// Get the latest version from https://github.com/sapics/ip-location-db/tree/master/dbip-city
	// Pattern: 34651136,34651647,BR,Sao Paulo,,Sao Paulo,,-23.5505,-46.6333,
	// City format: ip_range_start, ip_range_end, country_code, state1, state2, city, postcode, latitude, longitude, timezone
	fileCSV          = flag.String("file", "", "CSV file contain location database")
	ipVersion        = flag.String("ip-version", "4", "4 or 6 for IP CSV file version")
	databaseHost     = flag.String("host", "localhost", "Database host")
	databasePort     = flag.Int("port", 26257, "Database port")
	databaseUser     = flag.String("user", "root", "Database user")
	databasePassword = flag.String("password", "", "Database password")
	databaseBase     = flag.String("base", "corgi", "Database base")
)

func main() {
	flag.Parse()
	logger := log.New().With(context.TODO())

	if *ipVersion == "" {
		fmt.Println("please, specify ipv4 or ipv6 in ip-version flag")
		os.Exit(2)
	}

	if *ipVersion != "4" && *ipVersion != "6" {
		fmt.Println("you need pass ipv4 or ipv6 in ip-version flag")
		os.Exit(2)
	}

	// Get database connection.
	c := config.Config{}
	c.Database.Host = *databaseHost
	c.Database.Port = *databasePort
	c.Database.User = *databaseUser
	c.Database.Password = *databasePassword
	c.Database.Base = *databaseBase

	db := database.InitDatabase(logger, c)
	err := db.AutoMigrate(entity.LocationIPv4{}, entity.LocationIPv6{})
	if err != nil {
		fmt.Println("error to auto migrate: ", err.Error())
		os.Exit(2)
	}
	fmt.Println("File CSV: ", *fileCSV)

	if ok, _ := validFile(*fileCSV); !ok {
		fmt.Println("Hey fella, you need pass a valid location for CSV file")
		os.Exit(2)
	}

	file, err := os.Open(*fileCSV)
	if err != nil {
		fmt.Printf("error to open CSV file: %s\n", err.Error())
		os.Exit(2)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	if *ipVersion == "4" {
		var ips []entity.LocationIPv4
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				if db.Debug().Model(entity.LocationIPv4{}).Updates(&ips).RowsAffected == 0 {
					if err = db.Debug().Model(entity.LocationIPv4{}).Create(&ips).Error; err != nil {
						fmt.Println("error to create item:", err.Error())
						os.Exit(2)
					}
				}
				break
			}
			e, err := parseIPv4(text)
			if err != nil {
				fmt.Println("error: ", err.Error())
				os.Exit(2)
			}
			ips = append(ips, e)

			if len(ips) >= 5000 {
				if db.Debug().Model(entity.LocationIPv4{}).Updates(&ips).RowsAffected == 0 {
					if err = db.Debug().Model(entity.LocationIPv4{}).Create(&ips).Error; err != nil {
						fmt.Println("error to create item:", err.Error())
						os.Exit(2)
					}
				}
				ips = nil
			}
		}
	} else {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				break
			}
			e, err := parseIPv6(text)
			if err != nil {
				fmt.Println("error: ", err.Error())
				os.Exit(2)
			}
			if db.Debug().Model(entity.LocationIPv6{}).Where("range_start = ?", e.RangeStart).Updates(&e).RowsAffected == 0 {
				e.ID = uuid.New().String()
				if err = db.Debug().Model(entity.LocationIPv6{}).Create(&e).Error; err != nil {
					fmt.Println("error to create item:", err.Error())
					os.Exit(2)
				}
			}
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error to scan file: %s\n", err.Error())
		os.Exit(2)
	}
}

func validFile(path string) (ok bool, err error) {
	dir, err := os.Stat(path)
	if err != nil {
		return ok, fmt.Errorf("failed to open file, error: %w", err)
	}
	if dir.IsDir() {
		return ok, fmt.Errorf("%q is a directory", dir.Name())
	}
	return true, nil
}

func parseIPv4(text string) (e entity.LocationIPv4, err error) {
	line := strings.Split(text, ",")

	u, err := strconv.ParseUint(line[0], 10, 32)
	if err != nil {
		return
	}
	e.RangeStart = uint32(u)

	u, err = strconv.ParseUint(line[1], 10, 32)
	if err != nil {
		return
	}
	e.RangeEnd = uint32(u)

	e.Country = line[2]
	e.State = line[3]
	e.City = line[5]

	if i, err := strconv.ParseFloat(line[7], 64); err == nil {
		e.Latitude = i
	}

	if i, err := strconv.ParseFloat(line[8], 64); err == nil {
		e.Longitude = i
	}
	return
}

func parseIPv6(text string) (e entity.LocationIPv6, err error) {
	line := strings.Split(text, ",")

	e.RangeStart = line[0]
	e.RangeEnd = line[1]

	e.Country = line[2]
	e.State = line[3]
	e.City = line[5]

	if i, err := strconv.ParseFloat(line[7], 64); err == nil {
		e.Latitude = i
	}

	if i, err := strconv.ParseFloat(line[8], 64); err == nil {
		e.Longitude = i
	}
	return
}
