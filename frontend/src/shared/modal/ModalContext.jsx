import { createContext, useContext, useState } from "react"
import ModalRenderer from "./ModalRenderer"

const ModalContext = createContext()

export function ModalProvider({ children }) {
  const [modal, setModal] = useState(null)

  const confirm = (message) => {
    return new Promise((resolve) => {
      setModal({
        type: "confirm",
        message,
        resolve,
      })
    })
  }

  const success = (message) => {
    return new Promise((resolve) => {
      setModal({
        type: "success",
        message,
        resolve,
      })
    })
  }

  const error = (message) => {
    return new Promise((resolve) => {
      setModal({
        type: "error",
        message,
        resolve,
      })
    })
  }

  const loading = (message) => {
    setModal({
      type: "loading",
      message,
    })

    return () => setModal(null)
  }

  const close = (result) => {
    modal?.resolve?.(result)
    setModal(null)
  }

  return (
    <ModalContext.Provider value={{ confirm, success, error, loading }}>
      {children}
      {modal && <ModalRenderer modal={modal} close={close} />}
    </ModalContext.Provider>
  )
}

export function useModal() {
  return useContext(ModalContext)
}