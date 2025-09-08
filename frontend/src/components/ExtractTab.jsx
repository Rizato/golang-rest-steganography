import { useState } from 'react'
import apiService from '../services/api'

function ExtractTab() {
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
    
    if (!file) {
      setResult({ type: 'error', message: 'Please select an image file' })
      return
    }

    setLoading(true)
    setResult(null)
    setProgress('')

    try {
      const { message } = await apiService.extractMessage(
        file,
        (progressMessage) => setProgress(progressMessage)
      )

      setResult({ 
        type: 'success', 
        message: message
      })
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
      <h2>Extract Message</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Image file with hidden message:</label>
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
              id="extract-file-input"
            />
            <label htmlFor="extract-file-input" style={{ cursor: 'pointer' }}>
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
                  <small>Select an image that may contain a hidden message</small>
                </div>
              )}
            </label>
          </div>
        </div>

        <button type="submit" className="btn" disabled={loading}>
          {loading ? (progress || 'Processing...') : 'Extract Message'}
        </button>
      </form>

      {result && (
        <div className={`result ${result.type}`}>
          {result.type === 'success' && result.message ? (
            <div>
              <strong>Extracted Message:</strong>
              <div style={{ 
                marginTop: '10px', 
                padding: '15px', 
                background: 'white', 
                border: '1px solid #ddd', 
                borderRadius: '4px',
                fontFamily: 'monospace',
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-word'
              }}>
                {result.message}
              </div>
            </div>
          ) : (
            result.message
          )}
        </div>
      )}
    </div>
  )
}

export default ExtractTab