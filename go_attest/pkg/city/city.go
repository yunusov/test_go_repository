package city

import "fmt"

type City struct {
	Id         string `csv:"id"`
	Name       string `csv:"name"`
	Region     string `csv:"region"`
	District   string `csv:"district"`
	Population int    `csv:"population"`
	Foundation int    `csv:"foundation"`
}

func (c *City) GetId() string {
	return c.Id
}

func (c *City) GetName() string {
	return c.Name
}

func (c *City) GetRegion() string {
	return c.Region
}

func (c *City) GetDistrict() string {
	return c.District
}

func (c *City) GetPopulation() int {
	return c.Population
}

func (c *City) GetFoundation() int {
	return c.Foundation
}

func (c *City) UpdatePopulation(population int) {
	c.Population = population
}

func NewCity() (c *City) {
	return &City{}
}

func (c *City) ToString() string {
	return fmt.Sprintf("Id = %s, Name = %s, Region = %s, District = %s, "+
		"Population = %d, Foundation = %d\n", c.Id, c.Name, c.Region,
		c.District, c.Population, c.Foundation)
}
