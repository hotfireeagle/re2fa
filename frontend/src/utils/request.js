const SUCCESS_CODE = 1

export default function request(url, data={}, method="post") {
  const uppercaseMethod = method.toUpperCase()

  const requestParams = {
    method: uppercaseMethod,
    headers: {
      "Content-Type": "application/json"
    }
  }

  if (uppercaseMethod === "POST") {
    requestParams.body = JSON.stringify(data)
  } else if (uppercaseMethod === "GET") {
    const keys = Object.keys(data)
    const queryString = keys.map(key => `${key}=${data[key]}`).join("&")
    url = `${url}?${queryString}`
  } else {
    console.error("[todo]:implement other methods")
    return Promise.reject()
  }

  return fetch(url, requestParams)
    .then(response => response.json())
    .then(result => {
      if (result.code == SUCCESS_CODE) {
        return result.data
      } else {
        console.error(result.errorLog)
        alert(result.msg)
        return Promise.reject(result.msg)
      }
    })
}