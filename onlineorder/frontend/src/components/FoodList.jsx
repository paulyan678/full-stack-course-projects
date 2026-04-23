import { PlusOutlined, ShopOutlined } from "@ant-design/icons";
import { Button, Card, Empty, Select, Skeleton, message } from "antd";
import { useEffect, useState } from "react";
import { addItemToCart, getMenus, getRestaurants } from "../api/client";

function formatPrice(value) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD"
  }).format(Number(value));
}

function AddButton({ itemId, onAdded }) {
  const [loading, setLoading] = useState(false);

  async function add() {
    setLoading(true);
    try {
      await addItemToCart(itemId);
      message.success("Added to your cart");
      onAdded();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Button
      aria-label="Add item to cart"
      icon={<PlusOutlined />}
      loading={loading}
      onClick={add}
      shape="circle"
      type="primary"
    />
  );
}

export default function FoodList({ onCartChanged }) {
  const [restaurants, setRestaurants] = useState([]);
  const [restaurantId, setRestaurantId] = useState();
  const [items, setItems] = useState([]);
  const [loadingRestaurants, setLoadingRestaurants] = useState(true);
  const [loadingItems, setLoadingItems] = useState(false);

  useEffect(() => {
    let active = true;
    getRestaurants()
      .then((data) => {
        if (!active) return;
        setRestaurants(data);
        if (data.length) setRestaurantId(data[0].id);
      })
      .catch((error) => active && message.error(error.message))
      .finally(() => active && setLoadingRestaurants(false));
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    if (!restaurantId) return undefined;
    let active = true;
    setLoadingItems(true);
    getMenus(restaurantId)
      .then((data) => active && setItems(data))
      .catch((error) => active && message.error(error.message))
      .finally(() => active && setLoadingItems(false));
    return () => {
      active = false;
    };
  }, [restaurantId]);

  const selected = restaurants.find((restaurant) => restaurant.id === restaurantId);

  return (
    <section aria-labelledby="menu-heading" className="menu-section">
      <div className="menu-toolbar">
        <div>
          <span className="eyebrow">Tonight's picks</span>
          <h1 id="menu-heading">Find your next favorite</h1>
          <p>Freshly prepared by restaurants around the South Bay.</p>
        </div>
        <Select
          aria-label="Restaurant"
          loading={loadingRestaurants}
          onChange={setRestaurantId}
          options={restaurants.map((restaurant) => ({
            label: restaurant.name,
            value: restaurant.id
          }))}
          placeholder="Choose a restaurant"
          suffixIcon={<ShopOutlined />}
          value={restaurantId}
        />
      </div>

      {selected && (
        <div className="restaurant-meta">
          <strong>{selected.name}</strong>
          <span>{selected.address}</span>
          <span>{selected.phone}</span>
        </div>
      )}

      {loadingItems ? (
        <div className="food-grid" aria-label="Loading menu">
          {[1, 2, 3, 4, 5, 6].map((item) => (
            <Card key={item}><Skeleton active /></Card>
          ))}
        </div>
      ) : items.length ? (
        <div className="food-grid">
          {items.map((item) => (
            <Card
              className="food-card"
              cover={(
                <div className="food-image-wrap">
                  <img
                    alt=""
                    className="food-image"
                    loading="lazy"
                    onError={(event) => { event.currentTarget.hidden = true; }}
                    src={item.image_url}
                  />
                </div>
              )}
              key={item.id}
            >
              <div className="food-card-heading">
                <h2>{item.name}</h2>
                <AddButton itemId={item.id} onAdded={onCartChanged} />
              </div>
              <p>{item.description || "A house favorite, made fresh to order."}</p>
              <strong className="price">{formatPrice(item.price)}</strong>
            </Card>
          ))}
        </div>
      ) : (
        <Empty description="No dishes are available for this restaurant." />
      )}
    </section>
  );
}
