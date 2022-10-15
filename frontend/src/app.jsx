import { useState, useEffect } from "react"
import request from "./utils/request"
import { Form, Input, Button, message, notification, Select, Tag } from "antd"
import useFA from "@/hooks/useFA"
import styles from "./app.module.css"
import { useGet } from "@/hooks/useGet"

export default function App() {
  const [headerFormInstance] = Form.useForm()
  const [footerFormInstance] = Form.useForm()
  const [fa1Res, setFa1Res] = useState(null)
  const [fa2Res, setFa2Res] = useState(null)
  const [apiListRes, loadingApiListRes] = useGet("/api/apiList", [])

  useEffect(() => {
    const obj = apiListRes?.[0] || {}
    headerFormInstance.setFieldsValue({ api: obj.api })
  }, [apiListRes])

  useFA(fa1Res)
  useFA(fa2Res)

  const actionHandler = () => {
    const postData = headerFormInstance.getFieldsValue()
    const hide = message.loading({ content: "processing...", duration: 0 })
    return request(postData.api, postData).then(res => {
      setFa1Res({ ui: res[0], id: "fa1", title: res[0].title })
      setFa2Res({ ui: res[1], id: "fa2", title: res[1].title })
    }).finally(() => {
      hide()
    })
  }

  const matchFeedback = (id, result) => {
    const posMap = {
      fa1: "bottomLeft",
      fa2: "bottomRight",
    }
    const resultMap = new Map().set(true, "success").set(false, "error")
    notification[resultMap.get(result["result"])]({
      message: result["result"] ? "Match!" : "Unmatch!",
      description: `执行时间: ${result["time"]}us.`,
      placement: posMap[id],
    })
  }

  const matchHandler = () => {
    const text = footerFormInstance.getFieldValue("text")
    const headerFormObj = headerFormInstance.getFieldsValue()
    const postData = {
      ...(headerFormObj || {}),
      text,
    }
    if (!regexp) {
      message.error("Please input regexp")
      return
    }
    return request("/api/match", postData).then(res => {
      matchFeedback("fa1", res[0])
      matchFeedback("fa2", res[1])
    })
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.header}>
        <Form
          form={headerFormInstance}
          layout="inline"
          onFinish={actionHandler}
        >
          <Form.Item name="api">
            <Select style={{ width: 200 }} loading={loadingApiListRes}>
              {
                apiListRes.map(item => (
                  <Select.Option value={item.api} key={item.api}>
                    {item.name}
                  </Select.Option>
                ))
              }
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
              Action
            </Button>
          </Form.Item>
        </Form>
      </div>
      <div className={styles.contentArea}>
        <div className={`${styles.fa} ${styles.fa1}`}>
          <div id="fa1" className={styles.fawrapper} />
          <div className={styles.tagCls}>
            {
              fa1Res?.title ? (
                <Tag color="#108ee9" style={{ marginRight:0 }}>{fa1Res.title}</Tag>
              ) : null
            }
          </div>
        </div>
        <div className={`${styles.fa} ${styles.fa2}`}>
          <div id="fa2" className={styles.fawrapper} />
          <div className={styles.tagCls}>
            {
              fa2Res?.title ? (
                <Tag color="#108ee9" style={{ marginRight:0 }}>{fa2Res.title}</Tag>
              ) : null
            }
          </div>
        </div>
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