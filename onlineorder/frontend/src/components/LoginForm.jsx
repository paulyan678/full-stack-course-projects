import { LockOutlined, MailOutlined } from "@ant-design/icons";
import { Button, Form, Input, message } from "antd";
import { useState } from "react";
import { login } from "../api/client";

export default function LoginForm({ onSuccess }) {
  const [loading, setLoading] = useState(false);

  async function onFinish(values) {
    setLoading(true);
    try {
      await login(values);
      message.success("Welcome back");
      onSuccess(values.email.trim().toLowerCase());
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Form
      aria-label="Sign in"
      layout="vertical"
      name="login"
      onFinish={onFinish}
      requiredMark={false}
      size="large"
    >
      <Form.Item
        label="Email"
        name="email"
        rules={[
          { required: true, message: "Enter your email" },
          { type: "email", message: "Enter a valid email" }
        ]}
      >
        <Input
          autoComplete="email"
          maxLength={254}
          prefix={<MailOutlined />}
          placeholder="you@example.com"
        />
      </Form.Item>
      <Form.Item
        label="Password"
        name="password"
        rules={[{ required: true, message: "Enter your password" }]}
      >
        <Input.Password
          autoComplete="current-password"
          maxLength={72}
          prefix={<LockOutlined />}
          placeholder="Your password"
        />
      </Form.Item>
      <Button block htmlType="submit" loading={loading} type="primary">
        Sign in
      </Button>
    </Form>
  );
}
