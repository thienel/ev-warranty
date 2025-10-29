import { FILE_UPLOAD_CONFIG } from '@constants/common-constants'

export interface FileValidationResult {
  isValid: boolean
  errors: string[]
}

export interface FileUploadOptions {
  maxSize?: number
  allowedTypes?: string[]
  maxFiles?: number
}

/**
 * Validate a single file against upload constraints
 */
export function validateFile(file: File, options: FileUploadOptions = {}): FileValidationResult {
  const {
    maxSize = FILE_UPLOAD_CONFIG.MAX_FILE_SIZE,
    allowedTypes = FILE_UPLOAD_CONFIG.ALLOWED_TYPES,
  } = options

  const errors: string[] = []

  // Check file size
  if (file.size > maxSize) {
    const maxSizeMB = Math.round(maxSize / (1024 * 1024))
    errors.push(`File size must be less than ${maxSizeMB}MB`)
  }

  // Check file type
  if (!allowedTypes.includes(file.type as never)) {
    errors.push(`File type ${file.type} is not allowed`)
  }

  return {
    isValid: errors.length === 0,
    errors,
  }
}

/**
 * Validate multiple files against upload constraints
 */
export function validateFiles(
  files: FileList | File[],
  options: FileUploadOptions = {},
): FileValidationResult {
  const { maxFiles = FILE_UPLOAD_CONFIG.MAX_FILES } = options

  const fileArray = Array.from(files)
  const errors: string[] = []

  // Check number of files
  if (fileArray.length > maxFiles) {
    errors.push(`Maximum ${maxFiles} files allowed`)
  }

  // Validate each file
  fileArray.forEach((file, index) => {
    const validation = validateFile(file, options)
    if (!validation.isValid) {
      validation.errors.forEach((error) => {
        errors.push(`File ${index + 1}: ${error}`)
      })
    }
  })

  return {
    isValid: errors.length === 0,
    errors,
  }
}

/**
 * Format file size in human readable format
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Get file extension from filename
 */
export function getFileExtension(filename: string): string {
  return filename.slice(((filename.lastIndexOf('.') - 1) >>> 0) + 2)
}

/**
 * Generate a preview URL for an image file
 */
export function createFilePreview(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!file.type.startsWith('image/')) {
      reject(new Error('File is not an image'))
      return
    }

    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = () => reject(new Error('Failed to read file'))
    reader.readAsDataURL(file)
  })
}

/**
 * Convert FileList to Array
 */
export function fileListToArray(fileList: FileList): File[] {
  return Array.from(fileList)
}

/**
 * Check if file type is an image
 */
export function isImageFile(file: File): boolean {
  return file.type.startsWith('image/')
}

/**
 * Check if file type is a document
 */
export function isDocumentFile(file: File): boolean {
  const documentTypes = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'text/plain',
  ]
  return documentTypes.includes(file.type)
}
