use std::{collections::HashMap, f64::consts::PI};

#[derive(Serialize, Deserialize, Clone, Debug)]
#[serde(untagged)]
pub enum Pan {
    Rect { length: f64, width: f64 },
    Circle { diameter: f64 },
    Cupcake { diameter: f64, count: i32 },
    Bunndt { inner: f64, outer: f64 },
}
pub fn inventory() -> HashMap<String, Pan> {
    let timber_resources: HashMap<String, Pan> = [
        (
            "9x13".to_string(),
            Pan::Rect {
                length: 9.0,
                width: 13.0,
            },
        ),
        ("9\"".to_string(), Pan::Circle { diameter: 9.0 }),
    ]
    .iter()
    .cloned()
    .collect();
    return timber_resources;
}
impl Pan {
    pub fn area(&self) -> f64 {
        match *self {
            Pan::Rect { length, width } => length * width,
            Pan::Circle { diameter } => diameter * diameter * PI,
            Pan::Cupcake { diameter, count } => Pan::Circle { diameter }.area() * f64::from(count),
            Pan::Bunndt { inner, outer } => {
                Pan::Circle { diameter: outer }.area() - Pan::Circle { diameter: inner }.area()
            }
        }
    }
}

#[cfg(test)]
mod tests {

    use super::*;
    #[test]
    fn test_area() {
        let pan = Pan::Rect {
            length: 10.0,
            width: 20.0,
        };
        assert_eq!(pan.area(), 200.0);
        let pan2 = Pan::Cupcake {
            diameter: 2.0,
            count: 2,
        };
        assert_eq!(pan2.area(), f64::from(8) * PI);
        assert_eq!(inventory().get("9x13").unwrap().area(), 117.0)
    }
}
