import { useState } from "react"
import { useLogin } from "./useLogin"
import TextField from "../../shared/ui/TextField"
import PasswordField from "../../shared/ui/PasswordField"
import Button from "../../shared/ui/Button"

export default function LoginPage() {
  const { login, loading } = useLogin()

  const [form, setForm] = useState({
    email: "",
    password: "",
  })

  const handleChange = (e) => {
    const { name, value } = e.target

    setForm((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await login(form)
    } catch (err) {
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-sm bg-white p-6 rounded-xl shadow">
        <h1 className="text-2xl font-bold mb-6 text-center">
          Login
        </h1>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <TextField
            label="Email"
            name="email"
            placeholder="account@mail.com"
            value={form.email}
            onChange={handleChange}
          />

          <PasswordField
            label="Password"
            name="password"
            placeholder="••••••••"
            value={form.password}
            onChange={handleChange}
          />

          <Button
            type="submit"
            loading={loading}
            className="w-full"
          >
            Login
          </Button>
        </form>
      </div>
    </div>
  )
}