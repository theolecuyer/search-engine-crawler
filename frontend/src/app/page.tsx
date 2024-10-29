"use client";

import { useState } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";

export default function Home() {
  const [isLoading, setIsLoading] = useState(false); //State for loading status
  const [website, setWebsite] = useState(""); //To store the website being crawled
  //const [searchQuery, setSearchQuery] = useState(""); //To store the user query (Not yet implmented)
  const [results, setResults] = useState(""); //State to hold search results

  const handleSearch = async () => {
    setIsLoading(true);

    try {
      const response = await fetch(`/api/handler`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ website: website }), // Pass the query as website
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error("Error calling API:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-black">
      <h1 className="text-4xl font-bold text-white mb-4">SiteSearch</h1>
      {!isLoading ? (
        <div className="flex gap-2">
          <input
            type="text"
            className="border border-gray-300 rounded-md p-2 w-64 text-black"
            placeholder="Search..."
            value={website}
            onChange={(e) => setWebsite(e.target.value)}
          />
          <button
            className="bg-blue-500 text-white rounded px-4 h-full transition duration-300 hover:bg-blue-600"
            onClick={handleSearch}
          >
            Search
          </button>
        </div>
      ) : (
        <div className="flex items-center">
          <div className="loader"></div>
          <span className="text-white ml-2">Crawling {website}...</span>
        </div>
      )}
      <ScrollArea className="h-[200px] w-[350px] rounded-md border p-4">
        {results && <p className="text-white">{results}</p>}
      </ScrollArea>
      <style jsx>{`
        .loader {
          border: 4px solid rgba(255, 255, 255, 0.3);
          border-top: 4px solid #fff;
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
  );
}
