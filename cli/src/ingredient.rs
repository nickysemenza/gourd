use serde::{Deserialize, Serialize};
#[derive(Debug, PartialEq, Default, Deserialize, Serialize)]
pub struct Ingredient {
    pub name: String,
    pub grams: f64,

    pub amount: f64,
    pub unit: String,
    pub adjective: String,

    pub is_recipe: bool,

    pub optional: bool,
}

impl Ingredient {
    pub fn to_string(&self) -> String {
        let mut i = String::from(format!("{}g ", self.grams));
        if self.amount != 0.0 && self.unit != "" {
            i.push_str(&format!("({} {}) ", self.amount, self.unit));
        }
        i.push_str(&self.name);
        if self.is_recipe {
            i.push_str(&"*");
        }
        if self.adjective != "" {
            i.push_str(&format!(", {}", &self.adjective));
        }
        if self.optional {
            i.push_str(", optional")
        }
        return i;
    }
    #[allow(dead_code)]
    pub fn from_json(json: &str) -> Result<Self, &'static str> {
        let i: Ingredient = serde_json::from_str(json).unwrap();
        Ok(i)
    }
    #[allow(dead_code)]
    pub fn to_json(&self) -> Result<String, &'static str> {
        Ok(serde_json::to_string(self).unwrap())
    }
    #[allow(dead_code)]
    pub fn from_str(line: &str) -> Result<Self, &'static str> {
        // "120g (1 cup) flour, sifted, optional"
        // e.g.
        // Xg (Y unit) name, adj, optional
        // Xg name, adj, optional
        // Xg name, adj
        // Xg name
        let v: Vec<&str> = line.split(", ").collect();
        if v.len() > 3 || v.len() == 0 {
            return Err("invalid input");
        }
        let mut i = Ingredient::default();
        for x in v[1..].iter() {
            if x.eq(&"optional") {
                i.optional = true
            } else {
                i.adjective = x.to_string();
            }
        }
        let first = v[0];

        let au_start_bytes = first.find("(").unwrap_or(0) + 1;
        let au_end_bytes = first.find(")").unwrap_or(au_start_bytes);
        let au: Vec<&str> = first[au_start_bytes..au_end_bytes].split(" ").collect();
        if au.len() == 2 {
            i.amount = au[0].parse::<f64>().unwrap();
            i.unit = au[1].to_string();
        }
        let mut first2 = String::from(first);
        if au_start_bytes != au_end_bytes {
            first2.replace_range(au_start_bytes - 1..au_end_bytes + 2, "");
        }

        let mut grams_name = first2.splitn(2, ' ');
        let mut grams_str = grams_name.next().unwrap().to_string();
        grams_str.pop();
        i.grams = grams_str.parse::<f64>().unwrap();
        i.name = grams_name.next().unwrap().to_string();
        if i.name.ends_with("*") {
            i.name.truncate(i.name.len() - 1);
            i.is_recipe = true;
        }

        Ok(i)
    }
}

#[cfg(test)]
mod tests {
    use super::Ingredient;
    fn assert_back_and_forth(i: Ingredient, raw: &str) {
        assert_eq!(i.to_string(), raw, "parse");
        assert_eq!(Ingredient::from_str(raw).unwrap(), i, "to string");
    }

    #[test]
    fn base_case() {
        assert_back_and_forth(
            Ingredient {
                name: String::from("flour"),
                grams: 0.0,
                amount: 0.0,
                unit: String::from(""),
                adjective: String::from(""),
                optional: false,
                is_recipe: false,
            },
            "0g flour",
        );
        assert_back_and_forth(
            Ingredient {
                name: String::from("all purpose flour"),
                grams: 1.0,
                amount: 0.0,
                unit: String::from(""),
                adjective: String::from(""),
                optional: false,
                is_recipe: false,
            },
            "1g all purpose flour",
        );
        assert_back_and_forth(
            Ingredient {
                name: String::from("flour"),
                grams: 3.2,
                amount: 0.0,
                unit: String::from(""),
                adjective: String::from(""),
                optional: false,
                is_recipe: false,
            },
            "3.2g flour",
        );
    }
    #[test]
    fn all_values() {
        assert_back_and_forth(
            Ingredient {
                name: String::from("flour"),
                grams: 120.0,
                amount: 1.0,
                unit: "cup".to_string(),
                adjective: "sifted".to_string(),
                optional: true,
                is_recipe: false,
            },
            "120g (1 cup) flour, sifted, optional",
        );
        assert_back_and_forth(
            Ingredient {
                name: String::from("egg whites"),
                grams: 60.0,
                amount: 2.0,
                unit: "whites".to_string(),
                adjective: "room temperature".to_string(),
                optional: false,
                is_recipe: true,
            },
            "60g (2 whites) egg whites*, room temperature",
        );
    }
}
