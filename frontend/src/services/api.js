const getApiBaseUrl = () => {
  // Check for runtime config first (Docker)
  if (window.CONFIG && window.CONFIG.BACKEND_URL && window.CONFIG.BACKEND_URL !== 'BACKEND_URL_PLACEHOLDER') {
    return window.CONFIG.BACKEND_URL
  }
  // Fallback to build-time env var (development)
  return import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
}

const API_BASE_URL = getApiBaseUrl()

class ApiService {
  async uploadImage(file) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await fetch(`${API_BASE_URL}/api/v1/images`, {
      method: 'POST',
      body: formData
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Upload failed: ${errorText}`)
    }

    return await response.json()
  }

  async createEmbedJob(imageUuid, message) {
    const response = await fetch(`${API_BASE_URL}/api/v1/embed`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        'image-uuid': imageUuid,
        'message': message
      })
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Create embed job failed: ${errorText}`)
    }

    return await response.json()
  }

  async startEmbedJob(jobUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/embed/${jobUuid}/start`, {
      method: 'POST'
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Start embed job failed: ${errorText}`)
    }

    return await response.json()
  }

  async getEmbedJob(jobUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/embed/${jobUuid}`)

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Get embed job failed: ${errorText}`)
    }

    return await response.json()
  }

  async downloadImage(imageUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/images/${imageUuid}/download`)

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Download failed: ${errorText}`)
    }

    return await response.blob()
  }

  async createExtractJob(imageUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/extract`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        'image-uuid': imageUuid
      })
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Create extract job failed: ${errorText}`)
    }

    return await response.json()
  }

  async startExtractJob(jobUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/extract/${jobUuid}/start`, {
      method: 'POST'
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Start extract job failed: ${errorText}`)
    }

    return await response.json()
  }

  async getExtractJob(jobUuid) {
    const response = await fetch(`${API_BASE_URL}/api/v1/extract/${jobUuid}`)

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Get extract job failed: ${errorText}`)
    }

    return await response.json()
  }

  // High-level methods that combine the workflow
  async embedMessage(file, message, onProgress = () => {}) {
    try {
      onProgress('Uploading image...')
      const imageData = await this.uploadImage(file)
      
      onProgress('Creating embed job...')
      const jobData = await this.createEmbedJob(imageData.uuid, message)
      
      onProgress('Starting embed process...')
      await this.startEmbedJob(jobData.uuid)
      
      // Poll for completion
      onProgress('Processing...')
      let job = await this.getEmbedJob(jobData.uuid)
      
      while (job.status === 'pending' || job.status === 'running') {
        await new Promise(resolve => setTimeout(resolve, 1000))
        job = await this.getEmbedJob(jobData.uuid)
      }

      if (job.status === 'failed') {
        throw new Error(`Embed failed: ${job['status-message'] || 'Unknown error'}`)
      }

      onProgress('Downloading result...')
      const resultBlob = await this.downloadImage(job['embedded-uuid'])
      
      return {
        job,
        resultBlob,
        originalFilename: file.name
      }
    } catch (error) {
      throw error
    }
  }

  async extractMessage(file, onProgress = () => {}) {
    try {
      onProgress('Uploading image...')
      const imageData = await this.uploadImage(file)
      
      onProgress('Creating extract job...')
      const jobData = await this.createExtractJob(imageData.uuid)
      
      onProgress('Starting extract process...')
      await this.startExtractJob(jobData.uuid)
      
      // Poll for completion
      onProgress('Processing...')
      let job = await this.getExtractJob(jobData.uuid)
      
      while (job.status === 'pending' || job.status === 'running') {
        await new Promise(resolve => setTimeout(resolve, 1000))
        job = await this.getExtractJob(jobData.uuid)
      }

      if (job.status === 'failed') {
        throw new Error(`Extract failed: ${job['status-message'] || 'Unknown error'}`)
      }

      return {
        job,
        message: job.message || 'No hidden message found'
      }
    } catch (error) {
      throw error
    }
  }
}

export default new ApiService()