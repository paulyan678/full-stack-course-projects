INSERT INTO restaurants (name, address, image_url, phone)
SELECT 'Burger King', '773 N Mathilda Ave, Sunnyvale, CA 94085',
       'https://cdn.doordash.com/media/store/header/10171.png', '(408) 736-0101'
WHERE NOT EXISTS (SELECT 1 FROM restaurants WHERE name = 'Burger King');

INSERT INTO restaurants (name, address, image_url, phone)
SELECT 'SGD Tofu House', '3450 El Camino Real #105, Santa Clara, CA 95051',
       'https://cdn.doordash.com/media/store/header/1579.jpg', '(408) 261-3030'
WHERE NOT EXISTS (SELECT 1 FROM restaurants WHERE name = 'SGD Tofu House');

INSERT INTO restaurants (name, address, image_url, phone)
SELECT 'Fashion Wok', '163 S Murphy Ave, Sunnyvale, CA 94086',
       'https://cdn.doordash.com/media/store/header/273997.jpg', '(408) 739-8866'
WHERE NOT EXISTS (SELECT 1 FROM restaurants WHERE name = 'Fashion Wok');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Chicken Fries - 9 Pc', 4.89,
       'White meat chicken in a light, crispy breading seasoned with savory spices and herbs.',
       'https://cdn.doordash.com/media/photos/f439436f-c5ab-47af-bac4-7b73ab60a24b-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Chicken Fries - 9 Pc');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Whopper Meal', 10.59,
       'Flame-grilled beef with tomatoes, lettuce, mayonnaise, ketchup, pickles, and onions.',
       'https://cdn.doordash.com/media/photos/f878a689-618b-4c70-a00f-e7b1f320adc9-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Whopper Meal');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Impossible Whopper', 7.99,
       'A flame-grilled plant-based patty with classic Whopper toppings.',
       'https://cdn.doordash.com/media/photos/5c306a5f-fdd2-41d2-a660-9762aaa8eee8-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Impossible Whopper');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'HERSHEYS Sundae Pie', 3.09,
       'Crunchy chocolate crust, chocolate creme, topping, and chocolate chips.',
       'https://cdn.doordash.com/media/photos/80b1670d-e9c0-4886-a5b7-1ad48edd24ca-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'HERSHEYS Sundae Pie');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Original Chicken Sandwich', 6.09,
       'Lightly breaded chicken with shredded lettuce and creamy mayonnaise.',
       'https://cdn.doordash.com/media/photos/3e437f54-fa4e-4e9d-bf80-8a1e5b120f32-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Original Chicken Sandwich');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Classic OREO Shake', 3.99,
       'A creamy, hand-spun OREO shake.',
       'https://cdn.doordash.com/media/photos/c3ad483f-bad7-44f1-96af-4c3dcfc63c6d-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Burger King'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Classic OREO Shake');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Original Soft Tofu', 17.06,
       'Tofu stew with your choice of meat and mushrooms, served with rice and banchan.',
       'https://cdn.doordash.com/media/photos/b7055ca9-3caf-4d9d-9c99-04be1e36dbbf-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Original Soft Tofu');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Combination Soft Tofu', 17.06,
       'Tofu stew with beef, shrimp, and clams, served with rice and banchan.',
       'https://cdn.doordash.com/media/photos/37ad1974-1395-4e5c-86ff-fdf120cf8c58-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Combination Soft Tofu');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Seafood Soft Tofu', 17.06,
       'Tofu stew with mussels, shrimp, and clams, served with rice and banchan.',
       'https://cdn.doordash.com/media/photos/96bc8289-1950-4b4f-823d-12f33349a5fe-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Seafood Soft Tofu');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Seafood Pancake', 20.27,
       'Squid, clam, imitation crab, and grilled onions fried in a savory batter.',
       'https://cdn.doordash.com/media/photos/0a94b7e9-903d-49b7-937a-7940c8b56ad5-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Seafood Pancake');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Beef Short Ribs', 29.36,
       'Grilled beef short ribs served with rice and an assortment of banchan.',
       'https://cdn.doordash.com/media/photos/6340c369-2485-4d60-afcf-ca9068448d84-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Beef Short Ribs');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'BBQ Beef and Vegetables in Stoneware', 20.27,
       'Rice, barbecue beef, and vegetables served sizzling in stoneware.',
       'https://cdn.doordash.com/media/photos/9844dd4e-3c74-4942-8f90-2b3f4be25049-retina-large-jpeg'
FROM restaurants r
WHERE r.name = 'SGD Tofu House'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'BBQ Beef and Vegetables in Stoneware');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Stir Fried Pork with Pepper', 13.99,
       'Tender pork stir-fried with peppers. Medium spicy.',
       'https://cdn.doordash.com/media/photos/5b34852e-d253-461c-8be8-1bb0bc5e39be-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Stir Fried Pork with Pepper');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Eggplant with Minced Pork', 14.99,
       'Eggplant with minced pork, garlic, and cilantro.',
       'https://cdn.doordash.com/media/photos/bf70f262-0c55-41e1-89bc-84c061ae485f-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Eggplant with Minced Pork');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Poached Fish Fillets in Sour Soup', 17.99,
       'Poached fish fillets in a bright, mildly spicy sour broth.',
       'https://cdn.doordash.com/media/photos/1acf9c6b-189d-4583-a151-7ef522c283d9-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Poached Fish Fillets in Sour Soup');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Stir Fried Beef with Pepper', 16.99,
       'Tender beef stir-fried with fragrant peppers. Very spicy.',
       'https://cdn.doordash.com/media/photos/7f05859d-5e83-476d-a45a-73a3eb8a94e0-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Stir Fried Beef with Pepper');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Poached Sliced Beef in Hot Chili Oil', 17.99,
       'Sliced beef poached with aromatics in hot chili oil. Very spicy.',
       'https://cdn.doordash.com/media/photos/89ad8679-346e-41d8-b98f-3501fff4b277-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Poached Sliced Beef in Hot Chili Oil');

INSERT INTO menu_items (restaurant_id, name, price, description, image_url)
SELECT r.id, 'Fried Rice', 9.50,
       'Fried rice with broccoli, peas, carrots, bok choy, and egg.',
       'https://cdn.doordash.com/media/photos/ec06c431-9426-4971-a129-920440e1c9ce-retina-large.jpg'
FROM restaurants r
WHERE r.name = 'Fashion Wok'
  AND NOT EXISTS (SELECT 1 FROM menu_items m WHERE m.restaurant_id = r.id AND m.name = 'Fried Rice');
