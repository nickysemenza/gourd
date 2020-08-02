package scraper

import "time"

type Recipe struct {
	Context         string      `json:"@context"`
	Type            string      `json:"@type"`
	Name            string      `json:"name"`
	DatePublished   string      `json:"datePublished,omitempty"`
	DateModified    string      `json:"dateModified,omitempty"`
	Description     string      `json:"description"`
	Image           interface{} `json:"image"`
	AggregateRating struct {
		Type        string      `json:"@type"`
		RatingValue interface{} `json:"ratingValue"`
		RatingCount interface{} `json:"ratingCount"`
	} `json:"aggregateRating"`
	RecipeCategory     interface{} `json:"recipeCategory"`
	RecipeYield        string      `json:"recipeYield"`
	RecipeIngredient   []string    `json:"recipeIngredient"`
	RecipeInstructions []struct {
		Type string `json:"@type"`
		Text string `json:"text"`
	} `json:"recipeInstructions"`
	Author struct {
		Type        string   `json:"@type"`
		Name        string   `json:"name"`
		JobTitle    string   `json:"jobTitle"`
		SameAs      []string `json:"sameAs"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
	} `json:"author"`
	Keywords      string `json:"keywords"`
	TotalTime     string `json:"totalTime,omitempty"`
	RecipeCuisine string `json:"recipeCuisine,omitempty"`
	Nutrition     struct {
		Context               string      `json:"@context"`
		Type                  string      `json:"@type"`
		Calories              int         `json:"calories"`
		UnsaturatedFatContent string      `json:"unsaturatedFatContent"`
		CarbohydrateContent   string      `json:"carbohydrateContent"`
		CholesterolContent    interface{} `json:"cholesterolContent"`
		FatContent            string      `json:"fatContent"`
		FiberContent          string      `json:"fiberContent"`
		ProteinContent        string      `json:"proteinContent"`
		SaturatedFatContent   string      `json:"saturatedFatContent"`
		SodiumContent         string      `json:"sodiumContent"`
		SugarContent          string      `json:"sugarContent"`
		TransFatContent       string      `json:"transFatContent"`
	} `json:"nutrition,omitempty"`
	Video struct {
		Context      string    `json:"@context"`
		Type         string    `json:"@type"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		ThumbnailURL string    `json:"thumbnailUrl"`
		UploadDate   time.Time `json:"uploadDate"`
		Duration     string    `json:"duration"`
	} `json:"video,omitempty"`
}

// https://transform.tools/json-to-rust-serde
// [
//   {
//     "@context": "http://schema.org/",
//     "@type": "Recipe",
//     "name": "The Food Lab's Chocolate Chip Cookies Recipe",
//     "datePublished": "2013-12-19T14:00:00",
//     "dateModified": "2020-04-15T18:07:54",
//     "description": "Classic chocolate chip cookies level up with brown butter and chopped chocolate.",
//     "image": [
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit.jpg",
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit-300x225.jpg",
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit-750x563.jpg",
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit-1500x1125.jpg",
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit-625x469.jpg",
//       "https://www.seriouseats.com/recipes/images/2015/12/20131213-chocolate-chip-cookies-food-lab-55-edit-200x150.jpg"
//     ],
//     "aggregateRating": {
//       "@type": "AggregateRating",
//       "ratingValue": "4.5545454545455",
//       "ratingCount": "110"
//     },
//     "recipeCategory": [
//       "Chocolate",
//       "Cookies",
//       "Chocolate Chip Cookies",
//       "Christmas"
//     ],
//     "recipeYield": "Makes about 28 cookies",
//     "recipeIngredient": [
//       "8 ounces unsalted butter (2 sticks; 225g)",
//       "1 standard ice cube (about 2 tablespoons; 30mL frozen water)",
//       "10 ounces all-purpose flour (about 2 cups; 280g)",
//       "3/4 teaspoon (3g) baking soda ",
//       "2 teaspoons Diamond Crystal kosher salt or 1 teaspoon table salt (4g)",
//       "5 ounces granulated sugar (about 3/4 cup; 140g)",
//       "2 large eggs (100g)",
//       "2 teaspoons (10mL) vanilla extract",
//       "5 ounces dark brown sugar (about 1/2 tightly packed cup plus 2 tablespoons; 140g)",
//       "8 ounces (225g) semisweet chocolate, roughly chopped with a knife into 1/2- to 1/4-inch chunks",
//       "Coarse sea salt, for garnish"
//     ],
//     "recipeInstructions": [
//       {
//         "@type": "HowToStep",
//         "text": "Melt butter in a medium saucepan over medium-high heat. Cook, gently swirling pan constantly, until particles begin to turn golden brown and butter smells nutty, about 5 minutes. Remove from heat and continue swirling the pan until the butter is a rich brown, about 15 seconds longer. Transfer to a medium bowl, whisk in ice cube, transfer to refrigerator, and allow to cool completely, about 20 minutes, whisking occasionally. (Alternatively, whisk over an ice bath to hasten the process.)"
//       },
//       {
//         "@type": "HowToStep",
//         "text": "Meanwhile, whisk together flour, baking soda, and salt in a large bowl. Place granulated sugar, eggs, and vanilla extract in the bowl of a stand mixer fitted with the whisk attachment. Whisk on medium-high speed until mixture is pale brownish-yellow and falls off the whisk in thick ribbons when lifted, about 5 minutes."
//       },
//       {
//         "@type": "HowToStep",
//         "text": "Fit paddle attachment onto mixer. When brown butter mixture has cooled (it should be just starting to turn opaque again and firm around the edges), add brown sugar and cooled brown butter to egg mixture in stand mixer. Mix on medium speed to combine, about 15 seconds. Add flour mixture and mix on low speed until just barely combined, with some dry flour still remaining, about 15 seconds. Add chocolate and mix on low speed until dough comes together, about 15 seconds longer. Transfer to an airtight container and refrigerate dough at least overnight and up to 3 days."
//       },
//       {
//         "@type": "HowToStep",
//         "text": "When ready to bake, adjust oven racks to upper- and lower-middle positions and preheat oven to 325°F. Using a 1-ounce ice cream scoop or a spoon, place scoops of cookie dough onto a nonstick or parchment-lined baking sheet. Each ball should measure approximately 3 tablespoons in volume, and you should be able to fit 6 to 8 balls on each sheet. Tear each ball in half to reveal a rougher surface, then stick them back together with the rough sides facing outward. Transfer to oven and bake until golden brown around edges but still soft, 13 to 16 minutes, rotating pans back to front and top to bottom halfway through baking."
//       },
//       {
//         "@type": "HowToStep",
//         "text": "Remove baking sheets from oven. While cookies are still hot, sprinkle very lightly with coarse salt and gently press salt down to embed. Let cool for 2 minutes, then transfer cookies to a wire rack to cool completely. "
//       },
//       {
//         "@type": "HowToStep",
//         "text": "Repeat steps 3 through 5 for remaining cookie dough. Allow cookies to cool completely before storing in an airtight container, plastic bag, or cookie jar at room temperature for up to 5 days."
//       }
//     ],
//     "author": {
//       "@type": "Person",
//       "name": "J. Kenji López-Alt",
//       "jobTitle": "Chief Culinary Advisor",
//       "sameAs": [
//         "https://www.facebook.com/kenjilopezalt/",
//         "https://www.twitter.com/@kenjilopezalt"
//       ],
//       "description": "J. Kenji López-Alt is a stay-at-home dad who moonlights as the Chief Culinary Consultant of Serious Eats and the Chef/Partner of Wursthall, a German-inspired California beer hall near his home in San Mateo. His first book,  The Food Lab: Better Home Cooking Through Science (based on his Serious Eats column of the same name) is a New York Times best-seller, recipient of a James Beard Award, and was named Cookbook of the Year in 2015 by the International Association of Culinary Professionals. Kenji's next project is a children’s book called Every Night is Pizza Night, to be released in 2020, followed by another big cookbook in 2021.",
//       "url": "https://www.seriouseats.com/user/profile/kenjilopezalt"
//     },
//     "keywords": "baking, brown butter, chocolate, comfort, cookie, dessert"
//   },
//   {
//     "@context": "http://schema.org",
//     "@type": "Recipe",
//     "name": "Cornmeal-Blueberry Pancakes",
//     "description": "Pancakes are so easy to make that they encourage experimentation. Enter this cornmeal-blueberry variation, which feels like a weekend treat but is well suited for the weekday morning rush. Here, 1/2 cup of cornmeal is swapped out for a portion of the all-purpose flour, giving the pancake a wonderful texture. Make sure to dust the blueberries in flour before adding them to the batter; it will ensure even distribution of the fruit across the pancakes.",
//     "author": {
//       "@type": "Person",
//       "name": "Alison Roman"
//     },
//     "image": "https://static01.nyt.com/images/2016/06/15/dining/blueberrypancakes-2/blueberrypancakes-2-mediumThreeByTwo440.jpg",
//     "totalTime": "PT15M",
//     "recipeYield": "4 servings",
//     "recipeCuisine": "",
//     "recipeCategory": "",
//     "keywords": "",
//     "aggregateRating": {
//       "@type": "AggregateRating",
//       "ratingValue": 5,
//       "ratingCount": 402
//     },
//     "nutrition": {
//       "@context": "http://schema.org",
//       "@type": "NutritionInformation",
//       "calories": 499,
//       "unsaturatedFatContent": "5 grams",
//       "carbohydrateContent": "80 grams",
//       "cholesterolContent": null,
//       "fatContent": "13 grams",
//       "fiberContent": "3 grams",
//       "proteinContent": "15 grams",
//       "saturatedFatContent": "7 grams",
//       "sodiumContent": "1065 milligrams",
//       "sugarContent": "25 grams",
//       "transFatContent": "0 grams"
//     },
//     "recipeIngredient": [
//       "1/2 cup/80 grams cornmeal",
//       "1 1/2 cups/192 grams all-purpose flour",
//       "3 tablespoons/45 grams sugar",
//       "1 1/2 teaspoons/6 grams baking powder",
//       "1 1/2 teaspoons/9 grams baking soda",
//       "1 1/4 teaspoons kosher salt",
//       "2 1/2 cups buttermilk",
//       "2 large eggs",
//       "3 tablespoons/43 grams unsalted butter, melted",
//       "1 1/2 cups blueberries",
//       "Vegetable, canola or coconut oil for the skillet"
//     ],
//     "recipeInstructions": [
//       {
//         "@context": "http://schema.org",
//         "@type": "HowToStep",
//         "text": "Heat the oven to 325 degrees. Whisk cornmeal, flour, sugar, baking powder, baking soda and kosher salt together in a bowl. Using the whisk, make a well in the center. Pour the buttermilk into the well and crack eggs into buttermilk. Pour the melted butter into the mixture. Starting in the center, whisk everything together, moving toward the outside of the bowl, until all ingredients are incorporated. Do not overbeat. (Lumps are fine.) Coat your blueberries in a teaspoon of flour so that they don't sink, then stir them into the batter. The batter can be refrigerated for up to one hour."
//       },
//       {
//         "@context": "http://schema.org",
//         "@type": "HowToStep",
//         "text": "Heat a large nonstick griddle or skillet, preferably cast iron, over low heat for about 5 minutes. Add 1 tablespoon oil to the skillet. Turn heat up to medium–low and using a measuring cup, ladle 1/3 cup batter into the skillet. If you are using a large skillet or a griddle, repeat once or twice, taking care not to overcrowd the cooking surface."
//       },
//       {
//         "@context": "http://schema.org",
//         "@type": "HowToStep",
//         "text": "Flip pancakes after bubbles rise to surface and bottoms brown, after about 2 to 4 minutes. Cook until the other sides are lightly browned. Remove pancakes to a wire rack set inside a rimmed baking sheet, and keep in heated oven until all the batter is cooked and you are ready to serve."
//       }
//     ],
//     "video": {
//       "@context": "http://schema.org",
//       "@type": "VideoObject",
//       "name": "Cornmeal Blueberry Pancakes",
//       "description": "How to make cornmeal-blueberry buttermilk pancakes.",
//       "thumbnailUrl": "https://static01.nyt.com/images/2016/06/15/dining/blueberrypancakes-2/blueberrypancakes-2-mediumThreeByTwo440.jpg",
//       "uploadDate": "2016-06-10T20:36:39.000Z",
//       "duration": "PT1M51S"
//     }
//   }
// ]
