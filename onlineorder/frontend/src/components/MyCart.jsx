import { ShoppingOutlined } from "@ant-design/icons";
import { Badge, Button, Drawer, Empty, List, Spin, Typography, message } from "antd";
import { useEffect, useState } from "react";
import { checkout, getCart } from "../api/client";

const { Text } = Typography;

function money(value) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD"
  }).format(Number(value || 0));
}

export default function MyCart({ refreshKey }) {
  const [open, setOpen] = useState(false);
  const [cart, setCart] = useState();
  const [loading, setLoading] = useState(false);
  const [checkingOut, setCheckingOut] = useState(false);

  useEffect(() => {
    let active = true;
    setLoading(true);
    getCart()
      .then((data) => active && setCart(data))
      .catch((error) => active && message.error(error.message))
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [refreshKey]);

  async function onCheckout() {
    setCheckingOut(true);
    try {
      await checkout();
      setCart((current) => ({ ...current, total_price: 0, order_items: [] }));
      message.success("Order placed — thank you!");
    } catch (error) {
      message.error(error.message);
    } finally {
      setCheckingOut(false);
    }
  }

  const itemCount = cart?.order_items?.reduce((sum, item) => sum + item.quantity, 0) || 0;

  return (
    <>
      <Badge count={itemCount} overflowCount={99} size="small">
        <Button
          aria-label="Open shopping cart"
          className="cart-button"
          icon={<ShoppingOutlined />}
          onClick={() => setOpen(true)}
          shape="round"
        >
          Cart
        </Button>
      </Badge>
      <Drawer
        className="cart-drawer"
        footer={(
          <div className="cart-footer">
            <div>
              <span>Total</span>
              <Text strong>{money(cart?.total_price)}</Text>
            </div>
            <Button
              disabled={loading || !cart?.order_items?.length}
              loading={checkingOut}
              onClick={onCheckout}
              size="large"
              type="primary"
            >
              Checkout
            </Button>
          </div>
        )}
        onClose={() => setOpen(false)}
        open={open}
        title="Your order"
        width={480}
      >
        {loading ? (
          <div className="cart-loading"><Spin /></div>
        ) : cart?.order_items?.length ? (
          <List
            dataSource={cart.order_items}
            renderItem={(item) => (
              <List.Item key={item.order_item_id}>
                <List.Item.Meta
                  avatar={item.menu_item_image_url ? (
                    <img alt="" className="cart-thumbnail" src={item.menu_item_image_url} />
                  ) : null}
                  description={`${item.quantity} × ${money(item.price)}`}
                  title={item.menu_item_name}
                />
                <Text strong>{money(item.line_total)}</Text>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="Your cart is ready for something delicious." />
        )}
      </Drawer>
    </>
  );
}
