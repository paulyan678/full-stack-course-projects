import { LockOutlined, MailOutlined, UserOutlined } from "@ant-design/icons";
import { Button, Form, Input, message, Modal } from "antd";
import { useState } from "react";
import { signup } from "../api/client";

export default function SignupModal() {
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form] = Form.useForm();

  async function onFinish(values) {
    setSaving(true);
    try {
      await signup(values);
      message.success("Account created. You can sign in now.");
      form.resetFields();
      setOpen(false);
    } catch (error) {
      message.error(error.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <Button className="register-button" onClick={() => setOpen(true)}>
        Create account
      </Button>
      <Modal
        destroyOnClose
        footer={null}
        onCancel={() => setOpen(false)}
        open={open}
        title="Create your Lai Food account"
      >
        <Form
          form={form}
          layout="vertical"
          name="register"
          onFinish={onFinish}
          preserve={false}
          requiredMark={false}
        >
          <div className="name-fields">
            <Form.Item
              label="First name"
              name="first_name"
              rules={[{ required: true, message: "Enter your first name" }]}
            >
              <Input autoComplete="given-name" maxLength={80} prefix={<UserOutlined />} />
            </Form.Item>
            <Form.Item
              label="Last name"
              name="last_name"
              rules={[{ required: true, message: "Enter your last name" }]}
            >
              <Input autoComplete="family-name" maxLength={80} prefix={<UserOutlined />} />
            </Form.Item>
          </div>
          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: "Enter your email" },
              { type: "email", message: "Enter a valid email" }
            ]}
          >
            <Input autoComplete="email" maxLength={254} prefix={<MailOutlined />} />
          </Form.Item>
          <Form.Item
            extra="Use 8–72 characters."
            label="Password"
            name="password"
            rules={[
              { required: true, message: "Enter a password" },
              { min: 8, message: "Password must contain at least 8 characters" },
              { max: 72, message: "Password cannot exceed 72 characters" }
            ]}
          >
            <Input.Password autoComplete="new-password" prefix={<LockOutlined />} />
          </Form.Item>
          <Button block htmlType="submit" loading={saving} size="large" type="primary">
            Create account
          </Button>
        </Form>
      </Modal>
    </>
  );
}
