import { useState } from "react"
import request from "./utils/request"
import { Form, Input, Button, message, notification, Select } from "antd"
import useFANoEpsilon from "./hooks/useFANoEpsilon"
import styles from "./app.module.css"

export default function App() {
  const [headerFormInstance] = Form.useForm()
  const [footerFormInstance] = Form.useForm()
  const [nfaApiRes, setNfaApiRes] = useState(null)
  const [dfaApiRes, setDfaApiRes] = useState(null)

  useFANoEpsilon(nfaApiRes)
  useFANoEpsilon(dfaApiRes)

  const fetchFA = async () => {
    const postData = await headerFormInstance.validateFields()
    const urls = [
      { api: "/api/generateFA", id: "nfa", cb: setNfaApiRes, },
      { api: "/api/generateDFA", id: "dfa", cb: setDfaApiRes, },
    ]
    const hide = message.loading({ content: "processing...", duration: 0 })
    const requests = urls.map(obj => {
      let finalResult = null
      return request(obj.api, postData).then(result => {
        finalResult = result
      }).finally(() => {
        hide()
        obj.cb({ ui: finalResult, id: obj.id })
      })
    })
    return Promise.all(requests)
  }

  const matchFeedback = (id, result) => {
    const posMap = {
      nfa: "bottomLeft",
      dfa: "bottomRight",
    }
    const resultMap = new Map().set(true, "success").set(false, "error")
    notification[resultMap.get(result)]({
      message: result ? "Match!" : "Unmatch!",
      placement: posMap[id],
    })
  }

  const matchHandler = () => {
    const text = footerFormInstance.getFieldValue("text")
    const regexp = headerFormInstance.getFieldValue("regexp")
    const urls = [
      { api: "/api/nfaMatch", id: "nfa", },
      { api: "/api/dfaMatch", id: "dfa", },
    ]
    const fetchs = urls.map(obj => {
      let matchResult = false
      const postData = { regexp, text }
      return request(obj.api, postData).then(result => {
        matchResult = result
      }).finally(() => {
        matchFeedback(obj.id, matchResult)
      })
    })
    return Promise.all(fetchs)
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.header}>
        <Form
          form={headerFormInstance}
          layout="inline"
          onFinish={fetchFA}
        >
          <Form.Item name="mode">
            <Select style={{ width: 200 }}>
              <Select.Option>Generate DFA And NFA</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="regexp">
            <Input placeholder="Enter RegExp" />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
            >
              Generate DFA And NFA
            </Button>
          </Form.Item>
        </Form>
      </div>
      <div className={styles.contentArea}>
        <div id="nfa" className={`${styles.fa} ${styles.fa1}`} />
        <div id="dfa" className={`${styles.fa} ${styles.fa2}`} />
      </div>
      <div className={styles.footer}>
        <Form
          form={footerFormInstance}
          layout="inline"
          onFinish={matchHandler}
        >
          <Form.Item name="text">
            <Input placeholder="Enter Text" />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
            >
              Go To Match
            </Button>
          </Form.Item>
        </Form>
      </div>
    </div>
  )
}