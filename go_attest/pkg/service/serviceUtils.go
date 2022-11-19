package service

import (
	"fmt"
	city "go_attest/pkg/city"
	ut "go_attest/pkg/utils"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

const (
	csvFilename = "cities.csv"
	csvHeader   = "id,name,region,district,population,foundation\n"
)

func (s *Service) Init() {
	prepareCsvFile()
	in, err := os.Open(csvFilename)
	panicError(err)
	defer in.Close()

	cities := []*city.City{}
	err = gocsv.UnmarshalFile(in, &cities)
	panicError(err)
	for _, city := range cities {
		s.addCity(city)
	}
}

func (s *Service) SaveStore() {
	in, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	panicError(err)
	defer in.Close()
	store := []*city.City{}
	for _, v := range s.Store {
		store = append(store, v)
	}
	err = gocsv.MarshalFile(&store, in)
	panicError(err)
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareCsvFile() {
	in, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	panicError(err)
	defer in.Close()
	fi, err := in.Stat()
	panicError(err)
	bufSize := fi.Size()
	if bufSize == 0 {
		log.Printf("File %v is empty!", csvFilename)
		err = os.WriteFile(csvFilename, []byte(csvHeader), 0666)
		panicError(err)
		return
	}
	buf := make([]byte, bufSize)
	_, err = in.Read(buf)
	panicError(err)
	buffer := string(buf)
	bufHeader := buffer[:len(csvHeader)]

	if strings.Compare(bufHeader, csvHeader) != 0 {
		csvHeaderBuf := csvHeader + string(buf)
		err := os.WriteFile(csvFilename, []byte(csvHeaderBuf), 0666)
		panicError(err)
	}
}

func (s *Service) getCity(id string) (*city.City, error) {
	city := s.Store[id]
	if city == nil {
		return nil, fmt.Errorf("city is nil with ID=%v", id)
	}
	return city, nil
}

func (s *Service) addCity(city *city.City) {
	cityId := city.GetId()
	c := s.Store[cityId]
	if c == nil {
		s.Store[cityId] = city
	}
}

func (s *Service) deleteCity(cityId string) (string, error) {
	city, err := s.getCity(cityId)
	if err != nil {
		return "nil", err
	}
	delete(s.Store, cityId)
	return city.GetName(), nil
}

func (s *Service) updatePopulation(city *city.City, population int) {
	city.UpdatePopulation(population)
	s.Store[city.GetId()] = city
}

func (s *Service) getCities(paramValue string, paramName string) (result []*city.City) {
	for _, c := range s.Store {
		value := ""
		if paramName == "region" {
			value = c.GetRegion()
		} else if paramName == "district" {
			value = c.GetDistrict()
		}
		if strings.Compare(paramValue, value) == 0 {
			result = append(result, c)
		}
	}
	return
}

func (s *Service) getCitiesByParam(paramValue string, paramName string) (string, error) {
	cities := s.getCities(paramValue, paramName)
	response := ""
	for _, c := range cities {
		resp, err := ut.MarshalData(c)
		if err != nil {
			return "", err
		}
		response += resp
	}
	return response, nil
}

func (s *Service) getCitiesByRange(minValue int, maxValue int, paramName string) (string, error) {
	cities := s.getRangeCities(minValue, maxValue, paramName)
	response := ""
	for _, c := range cities {
		resp, err := ut.MarshalData(c)
		if err != nil {
			return "", err
		}
		response += resp
	}
	return response, nil
}

func (s *Service) getRangeCities(minValue int, maxValue int, paramName string) (result []*city.City) {
	for _, c := range s.Store {
		value := 0
		if paramName == "population" {
			value = c.GetPopulation()
		} else if paramName == "foundation" {
			value = c.GetFoundation()
		}
		if value >= minValue && value <= maxValue {
			result = append(result, c)
		}
	}
	return
}
