import { useState } from 'react'
import apiService from '../services/api'

function EmbedTab() {
  const [message, setMessage] = useState('')
  const [file, setFile] = useState(null)
  const [result, setResult] = useState(null)
  const [loading, setLoading] = useState(false)
  const [progress, setProgress] = useState('')

  const handleFileChange = (e) => {
    setFile(e.target.files[0])
  }

  const handleDrop = (e) => {
    e.preventDefault()
    const droppedFiles = e.dataTransfer.files
    if (droppedFiles.length > 0) {
      setFile(droppedFiles[0])
    }
  }

  const handleDragOver = (e) => {
    e.preventDefault()
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    if (!message.trim()) {
      setResult({ type: 'error', message: 'Please enter a message to embed' })
      return
    }
    
    if (!file) {
      setResult({ type: 'error', message: 'Please select an image file' })
      return
    }

    setLoading(true)
    setResult(null)
    setProgress('')

    try {
      const { resultBlob, originalFilename } = await apiService.embedMessage(
        file, 
        message, 
        (progressMessage) => setProgress(progressMessage)
      )

      // Download the result
      const url = window.URL.createObjectURL(resultBlob)
      const a = document.createElement('a')
      a.href = url
      a.download = `embedded_${originalFilename}`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
      
      setResult({ type: 'success', message: 'Message embedded successfully! Download started.' })
      setProgress('')
    } catch (error) {
      setResult({ type: 'error', message: error.message })
      setProgress('')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="tab-content">
      <h2>Embed Message</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="message">Message to embed:</label>
          <textarea
            id="message"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Enter your secret message here..."
            rows={4}
          />
        </div>
        
        <div className="form-group">
          <label>Image file:</label>
          <div 
            className="file-upload"
            onDrop={handleDrop}
            onDragOver={handleDragOver}
          >
            <input
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              style={{ display: 'none' }}
              id="file-input"
            />
            <label htmlFor="file-input" style={{ cursor: 'pointer' }}>
              {file ? (
                <div>
                  <strong>Selected:</strong> {file.name}
                  <br />
                  <small>Click to change or drag & drop a new file</small>
                </div>
              ) : (
                <div>
                  Click to select an image file or drag & drop here
                  <br />
                  <small>Supported formats: PNG, JPEG, GIF</small>
                </div>
              )}
            </label>
          </div>
        </div>

        <button type="submit" className="btn" disabled={loading}>
          {loading ? (progress || 'Processing...') : 'Embed Message'}
        </button>
      </form>

      {result && (
        <div className={`result ${result.type}`}>
          {result.message}
        </div>
      )}
    </div>
  )
}

export default EmbedTab