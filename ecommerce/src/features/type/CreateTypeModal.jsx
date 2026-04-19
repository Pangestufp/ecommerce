import { useState } from "react"
import TextField from "../../shared/ui/TextField"
import Button from "../../shared/ui/Button"

export default function CreateTypeModal({ onSubmit, onClose, loading }) {
  const [form, setForm] = useState({
    type_code: "",
    type_name: "",
    type_desc: "",
  })

  const handleChange = (e) => {
    const { name, value } = e.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    await onSubmit(form)
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[360px] shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-4">Buat Tipe</h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <TextField
            label="Kode Tipe"
            name="type_code"
            placeholder="Masukkan kode tipe"
            value={form.type_code}
            onChange={handleChange}
          />
          <TextField
            label="Nama Tipe"
            name="type_name"
            placeholder="Masukkan nama tipe"
            value={form.type_name}
            onChange={handleChange}
          />
          <TextField
            label="Deskripsi"
            name="type_desc"
            placeholder="Masukkan deskripsi tipe"
            value={form.type_desc}
            onChange={handleChange}
          />

          <div className="flex justify-end gap-2 mt-2">
            <Button variant="secondary" onClick={onClose} disabled={loading}>
              Batal
            </Button>
            <Button type="submit" loading={loading}>
              Simpan
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}