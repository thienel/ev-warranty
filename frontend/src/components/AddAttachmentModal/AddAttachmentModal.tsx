import React, { useState } from 'react'
import { Modal, Upload, Button, Form, message } from 'antd'
import { InboxOutlined } from '@ant-design/icons'
import type { UploadFile, UploadProps } from 'antd'
import { claimAttachments as claimAttachmentsApi } from '@/services/index'

const { Dragger } = Upload

interface AddAttachmentModalProps {
  visible: boolean
  claimId: string
  onCancel: () => void
  onSuccess: () => void
}

const AddAttachmentModal: React.FC<AddAttachmentModalProps> = ({
  visible,
  claimId,
  onCancel,
  onSuccess,
}) => {
  const [form] = Form.useForm()
  const [fileList, setFileList] = useState<UploadFile[]>([])
  const [uploading, setUploading] = useState(false)

  const uploadProps: UploadProps = {
    name: 'files',
    multiple: true,
    fileList,
    accept: 'image/*,video/*,.pdf,.doc,.docx,.txt', // Accept images, videos, PDFs, and documents
    beforeUpload: (file) => {
      // Debug file information
      console.log('File details:', {
        name: file.name,
        size: file.size,
        type: file.type,
        sizeInMB: (file.size / (1024 * 1024)).toFixed(2),
      })

      // Check file type first
      const isVideo =
        file.type.startsWith('video/') || /\.(mp4|avi|mov|wmv|flv|webm|mkv|3gp)$/i.test(file.name)
      const isImage =
        file.type.startsWith('image/') || /\.(jpg|jpeg|png|gif|bmp|webp|svg)$/i.test(file.name)
      const isDocument =
        file.type.includes('pdf') ||
        file.type.includes('document') ||
        file.type.includes('text') ||
        /\.(pdf|doc|docx|txt|rtf)$/i.test(file.name)

      // Validate file type
      if (!isVideo && !isImage && !isDocument) {
        message.error(
          `File "${file.name}" is not a supported format. Please upload images, videos, or documents only.`,
        )
        return false
      }

      // More generous size limits - only for frontend validation
      let maxSize: number
      let sizeLimitText: string

      if (isVideo) {
        maxSize = 50 * 1024 * 1024 // 50MB for videos
        sizeLimitText = '50MB'
      } else if (isImage) {
        maxSize = 10 * 1024 * 1024 // 10MB for images
        sizeLimitText = '10MB'
      } else {
        maxSize = 10 * 1024 * 1024 // 10MB for documents
        sizeLimitText = '10MB'
      }

      if (file.size > maxSize) {
        const fileSizeMB = (file.size / (1024 * 1024)).toFixed(2)
        message.error(
          `File "${file.name}" (${fileSizeMB}MB) is too large! Maximum size for ${isVideo ? 'videos' : isImage ? 'images' : 'documents'} is ${sizeLimitText}.`,
        )
        return false
      }

      return false // Prevent auto upload
    },
    onChange: ({ fileList: newFileList }) => {
      setFileList(newFileList)
    },
    onDrop: (e) => {
      console.log('Dropped files', e.dataTransfer.files)
    },
  }

  const handleUpload = async () => {
    if (fileList.length === 0) {
      message.error('Please select at least one file')
      return
    }

    try {
      setUploading(true)

      // Convert UploadFile[] to FileList
      const files = fileList.map((file) => file.originFileObj).filter(Boolean) as File[]

      // Debug files being uploaded
      console.log(
        'Files being uploaded:',
        files.map((f) => ({
          name: f.name,
          size: f.size,
          type: f.type,
          sizeInMB: (f.size / (1024 * 1024)).toFixed(2),
        })),
      )

      const fileListObj = new DataTransfer()
      files.forEach((file) => fileListObj.items.add(file))

      await claimAttachmentsApi.upload(claimId, fileListObj.files)

      message.success('Attachments uploaded successfully')
      onSuccess()
      handleReset()
    } catch (error: unknown) {
      console.error('Upload failed:', error)

      // Handle specific error types
      const axiosError = error as {
        response?: { status?: number; data?: unknown }
        message?: string
      }

      if (axiosError?.response?.status === 413) {
        message.error(
          'Server rejected: File too large! Even though frontend allows it, the server has a smaller limit. Try reducing file size further.',
        )
      } else if (axiosError?.response?.status === 404) {
        message.error('Upload endpoint not found. Please check your connection and try again.')
      } else if (axiosError?.response?.status === 400) {
        const errorData = axiosError?.response?.data
        console.log('Server error details:', errorData)
        message.error('Invalid file format or request. Please check your files and try again.')
      } else if (axiosError?.message?.includes('Network Error')) {
        message.error('Network error. Please check your connection and try again.')
      } else if (axiosError?.message?.includes('timeout')) {
        message.error('Upload timeout. Large files may take longer. Please try again.')
      } else {
        console.log('Full error object:', axiosError)
        message.error(`Upload failed: ${axiosError?.message || 'Unknown error'}. Please try again.`)
      }
    } finally {
      setUploading(false)
    }
  }

  const handleReset = () => {
    form.resetFields()
    setFileList([])
  }

  const handleCancel = () => {
    handleReset()
    onCancel()
  }

  return (
    <Modal
      title="Add Attachments"
      open={visible}
      onCancel={handleCancel}
      footer={[
        <Button key="cancel" onClick={handleCancel}>
          Cancel
        </Button>,
        <Button
          key="upload"
          type="primary"
          loading={uploading}
          onClick={handleUpload}
          disabled={fileList.length === 0}
        >
          Upload {fileList.length > 0 && `(${fileList.length} files)`}
        </Button>,
      ]}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="files"
          label="Select Files"
          rules={[{ required: true, message: 'Please select at least one file' }]}
        >
          <Dragger {...uploadProps}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">Click or drag files to this area to upload</p>
            <p className="ant-upload-hint">
              Support for single or bulk upload. Supported formats: Images, Videos, PDFs, Documents
              <br />
              Size limits: Videos up to 50MB, Images and Documents up to 10MB each
            </p>
          </Dragger>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default AddAttachmentModal
