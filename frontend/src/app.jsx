import { useState } from "react"
import request from "./utils/request"
import useFA from "./hooks/useFA"
import useFANoEpsilon from "./hooks/useFANoEpsilon"

export default function App() {
  const [regexp, setRegexp] = useState("")
  const [apiRes, setApiRes] = useState(null)
  const [faNoEpsilonApiRes, setFaNoEpsilonApiRes] = useState(null)

  useFA(apiRes)
  useFANoEpsilon(faNoEpsilonApiRes)

  const regexpChangeInputHandler = event => {
    setRegexp(event.target.value)
  }

  const fetchFA = (event, api, updater) => {
    event.preventDefault()
    let finalResult = null
    return request(api, { regexp }).then(result => {
      finalResult = result
    }).finally(() => {
      updater(finalResult)
    })
  }

  return (
    <div className="flex flex-col container bg-gray-50 w-screen max-w-none h-screen">
      <form className="w-screen p-5 bg-blue-600 bg-opacity-60">
        <div className="text-center">
          <input
            placeholder="Enter regexp"
            value={regexp}
            onChange={regexpChangeInputHandler}
            className="border w-96 max-w-none p-2 px-3 rounded-lg mr-6 border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
          />
          <button
            onClick={event => fetchFA(event, "/api/generateFA", setApiRes)}
            className="bg-blue-600 p-2 px-3 rounded-lg mr-6 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white"
          >
            Generate FA
          </button>
          <button
            onClick={event => fetchFA(event, "/api/generateFANoEpsilon", setFaNoEpsilonApiRes)}
            className="bg-blue-600 p-2 px-3 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white"
          >
            Generate FA But No Epsilon
          </button>
        </div>
      </form>
      <div id="fa" className="flex-grow bg-white" />
    </div>
  )
}