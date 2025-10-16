import React from 'react'
import { Link } from 'react-router-dom'

function Navbar() {
  return (
    <nav className="bg-blue-600 text-white shadow-lg">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="text-xl font-bold">
            Smart Task Orchestrator
          </Link>
          <div className="space-x-4">
            <Link 
              to="/" 
              className="hover:bg-blue-700 px-3 py-2 rounded transition-colors"
            >
              Dashboard
            </Link>
            <Link 
              to="/create" 
              className="hover:bg-blue-700 px-3 py-2 rounded transition-colors"
            >
              Create Job
            </Link>
          </div>
        </div>
      </div>
    </nav>
  )
}

export default Navbar