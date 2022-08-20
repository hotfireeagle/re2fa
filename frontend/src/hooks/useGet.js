import { useEffect, useState } from "react"
import request from "@/utils/request"

export const useGet = (url, defaultVal={}) => {
  const [apiRes, setApiRes] = useState(defaultVal)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    setLoading(true)
    request(url, {}, "get").then(res => {
      setApiRes(res)
    }).finally(() => {
      setLoading(false)
    })
  }, [url])

  return [apiRes, loading]
}
