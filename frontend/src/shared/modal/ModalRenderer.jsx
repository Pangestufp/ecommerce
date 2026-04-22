import Button from "../ui/Button"

export default function ModalRenderer({ modal, close }) {
  const { type, message } = modal

  const color = {
    confirm: "text-gray-800",
    success: "text-green-600",
    error: "text-red-600",
    loading: "text-blue-600",
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">

      <div className="bg-white rounded-lg p-6 w-[340px] shadow-lg text-center">

        <p className={`mb-6 font-medium ${color[type]}`}>
          {message}
        </p>

        {type === "confirm" && (
          <div className="flex justify-center gap-3">
            <Button
              variant="secondary"
              onClick={() => close(false)}
            >
              Cancel
            </Button>

            <Button
              variant="danger"
              onClick={() => close(true)}
            >
              Confirm
            </Button>
          </div>
        )}

        {(type === "success" || type === "error") && (
          <Button onClick={() => close(true)}>
            OK
          </Button>
        )}

        {type === "loading" && (
          <div className="flex justify-center">
            <div className="animate-spin rounded-full h-6 w-6 border-4 border-blue-500 border-t-transparent"></div>
          </div>
        )}

      </div>

    </div>
  )
}