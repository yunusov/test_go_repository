package service

import (
	city "go_attest/pkg/city"
	ut "go_attest/pkg/utils"
	"log"
	"net/http"
)

type Service struct {
	Store map[string]*city.City
}

func NewService() *Service {
	return &Service{make(map[string]*city.City)}
}

func (s *Service) GetCityById(w http.ResponseWriter, r *http.Request) {
	/*
		получение информации о городе по его id;
	*/
	ut.LogRequest("GetCityById", r)
	defer r.Body.Close()

	cityId := ut.GetRequestParam(r, "cityid")
	log.Printf("cityId = %v", cityId)
	c, err := s.getCity(cityId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	response, err := ut.MarshalData(c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		добавление новой записи в список городов;
	*/
	ut.LogRequest("Create", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		content, err := ut.GetContent(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		city := city.NewCity()
		if err := ut.UnMarshalData(content, &city, w); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		s.addCity(city)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(city.GetId()))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	/*
		удаление информации о городе по указанному id;
	*/
	ut.LogRequest("Delete", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		cityId := dat["target_id"].(string)
		userName, err := s.deleteCity(cityId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userName))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) UpdatePopulationById(w http.ResponseWriter, r *http.Request) {
	/*
		обновление информации о численности населения города по указанному id;
	*/
	ut.LogRequest("UpdatePopulationById", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		cityid := ut.GetRequestParam(r, "cityid")
		log.Printf("cityid = %v", cityid)
		city, err := s.getCity(cityid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		s.updatePopulation(city, int(dat["new population"].(float64)))
		response := "Param 'population' was changed successfully!"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetByRegion(w http.ResponseWriter, r *http.Request) {
	/*
		получение списка городов по указанному региону;
	*/
	ut.LogRequest("GetByRegion", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		paramValue := dat["region"].(string)
		response, err := s.getCitiesByParam(paramValue, "region")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetByDistrict(w http.ResponseWriter, r *http.Request) {
	/*
		получение списка городов по указанному округу;
	*/
	ut.LogRequest("GetByDistrict", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		paramValue := dat["district"].(string)
		response, err := s.getCitiesByParam(paramValue, "district")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetByPopulationRange(w http.ResponseWriter, r *http.Request) {
	/*
		получения списка городов по указанному диапазону численности населения;
	*/
	ut.LogRequest("GetByPopulationRange", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		minValue := int(dat["min"].(float64))
		maxValue := int(dat["max"].(float64))
		response, err := s.getCitiesByRange(minValue, maxValue, "population")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetByFoundationRange(w http.ResponseWriter, r *http.Request) {
	/*
		получения списка городов по указанному диапазону года основания.
	*/
	ut.LogRequest("GetByPopulationRange", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		minValue := int(dat["min"].(float64))
		maxValue := int(dat["max"].(float64))
		response, err := s.getCitiesByRange(minValue, maxValue, "foundation")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
