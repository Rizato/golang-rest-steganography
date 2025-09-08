import { useState } from 'react'
import './App.css'
import EmbedTab from './components/EmbedTab'
import ExtractTab from './components/ExtractTab'

function App() {
  const [activeTab, setActiveTab] = useState('embed')

  return (
    <div className="app">
      <header className="app-header">
        <h1>Steganography Tool</h1>
        <nav className="tabs">
          <button 
            className={activeTab === 'embed' ? 'active' : ''}
            onClick={() => setActiveTab('embed')}
          >
            Embed
          </button>
          <button 
            className={activeTab === 'extract' ? 'active' : ''}
            onClick={() => setActiveTab('extract')}
          >
            Extract
          </button>
        </nav>
      </header>
      <main className="app-main">
        {activeTab === 'embed' && <EmbedTab />}
        {activeTab === 'extract' && <ExtractTab />}
      </main>
    </div>
  )
}

export default App