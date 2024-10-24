package port

type CategoryEntity struct {
	Name     string
	MccCodes []string
	Order    int
}

/*
	TODO:
	In a perfect world, the sections below could be register in a `categories` table/document/data(source|storage)
	that should be related to an `mcc_codes`,  ensuring  flexibility  when creating new categories. For simplicity
	in  developing  the  core  of  the  `services`  and  `domains`,  I  have  left  them hard-coded here, to later
	extract into a more appropriate approach.
*/

const (
	CATEGORY_FOOD_NAME string = "FOOD"
	CATEGORY_MEAL_NAME string = "MEAL"
	CATEGORY_CASH_NAME string = "CASH"
)

var (
	CategoryFood = CategoryEntity{Name: CATEGORY_FOOD_NAME, MccCodes: []string{"5411", "5412"}, Order: 1}
	CategoryMeal = CategoryEntity{Name: CATEGORY_MEAL_NAME, MccCodes: []string{"5811", "5812"}, Order: 2}
	CategoryCash = CategoryEntity{Name: CATEGORY_CASH_NAME, Order: 3}
)

var Categories = map[string]CategoryEntity{
	CATEGORY_FOOD_NAME: CategoryFood,
	CATEGORY_MEAL_NAME: CategoryMeal,
	CATEGORY_CASH_NAME: CategoryCash,
}
