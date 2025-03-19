"use client"

import { useState } from "react"
import { ScrollArea } from "@/components/ui/scroll-area"
import axios from "axios"
import Image from "next/image"

export default function Home() {
  const [isLoading, setIsLoading] = useState(false)
  const [searchTerm, setsearchTerm] = useState("")
  const [results, setResults] = useState([])
  const [hasSearched, setHasSearched] = useState(false)

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
  const apiKey = process.env.NEXT_PUBLIC_GO_API_KEY || ""
  const handleSearch = async () => {
    setIsLoading(true)
    console.log("handleSearch called")
    try {
      const response = await axios.post(
        apiUrl,
        { search_term: searchTerm },
        {
          headers: {
            Authorization: apiKey,
          },
        },
      )
      const data = response.data || []
      setResults(data)
    } catch (error) {
      console.error("Error calling API:", error)
      setResults([])
    } finally {
      setIsLoading(false)
      setHasSearched(true)
    }
  }

  return (
    <div className="flex flex-col min-h-screen bg-[#1a1a1a]">
      {!hasSearched ? (
        //layout before search
        <div className="flex flex-col items-center justify-center min-h-screen">
          <div className="flex flex-col items-center mb-8">
            <Image 
              src="/USFlogo.png"
              alt="USF Logo" 
              width={100} 
              height={100}
              className="mb-4"
            />
            <h1 className="text-4xl font-bold text-white">USF Crawled Web Search</h1>
          </div>
          {!isLoading ? (
            <div className="flex gap-2">
              <input
                type="text"
                className="border border-gray-700 bg-[#222222] rounded-md p-2 w-64 text-white"
                placeholder="Search courses..."
                value={searchTerm}
                onChange={(e) => setsearchTerm(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
              <button
                className="bg-[#00543C] text-white rounded px-4 py-2 flex items-center transition duration-300 hover:bg-[#006b4f]"
                onClick={handleSearch}
              >
                Search
              </button>
            </div>
          ) : (
            <div className="flex items-center mt-4">
              <div className="loader"></div>
              <span className="text-white ml-2">Searching {searchTerm}...</span>
            </div>
          )}
        </div>
      ) : (
        //Layout after search is done
        <div className="flex flex-col w-full">
          {/* Header with logo and search bar */}
          <div className="flex items-center p-4 bg-[#00543C] border-b border-[#006b4f]">
            <div className="flex items-center mr-8">
              <Image 
                src="/USFlogo.png"
                alt="USF Logo" 
                width={40} 
                height={40}
                className="mr-3"
              />
              <h1 className="text-xl font-bold text-white whitespace-nowrap">USF Crawled Web Search</h1>
            </div>
            {!isLoading ? (
              <div className="flex gap-2 flex-1 max-w-md">
                <input
                  type="text"
                  className="border border-gray-700 bg-[#222222] rounded-md p-2 flex-1 text-white"
                  placeholder="Search courses..."
                  value={searchTerm}
                  onChange={(e) => setsearchTerm(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSearch()}
                />
                <button
                  className="bg-[#1a1a1a] text-white rounded px-4 py-2 flex items-center transition duration-300 hover:bg-[#333333]"
                  onClick={handleSearch}
                >
                  Search
                </button>
              </div>
            ) : (
              <div className="flex items-center">
                <div className="loader"></div>
                <span className="text-white ml-2">Searching {searchTerm}...</span>
              </div>
            )}
          </div>

          {/* Results area - full width, no border */}
          <div className="flex-1">
            <ScrollArea className="h-[calc(100vh-80px)] w-full">
              <div className="p-4 max-w-5xl mx-auto">
                {results.length > 0 ? (
                  results.map((url, index) => (
                    <p key={index} className="py-2">
                      <a
                        href={url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-[#e6b800] hover:text-[#ffcc00] underline"
                      >
                        {url}
                      </a>
                    </p>
                  ))
                ) : (
                  <p className="text-white">No results found</p>
                )}
              </div>
            </ScrollArea>
          </div>
        </div>
      )}
      <style jsx>{`
        .loader {
          border: 4px solid rgba(255, 255, 255, 0.3);
          border-top: 4px solid #e6b800;
          border-radius: 50%;
          width: 20px;
          height: 20px;
          animation: spin 1s linear infinite;
        }

        @keyframes spin {
          0% {
            transform: rotate(0deg);
          }
          100% {
            transform: rotate(360deg);
          }
        }
      `}</style>
    </div>
  )
}
