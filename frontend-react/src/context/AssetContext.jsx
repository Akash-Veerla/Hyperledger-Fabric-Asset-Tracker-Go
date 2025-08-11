import { createContext, useState, useEffect, useContext } from 'react';

const AssetContext = createContext();

export const useAssets = () => {
  return useContext(AssetContext);
};

export const AssetProvider = ({ children }) => {
  const [assets, setAssets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const API_URL = 'http://localhost:8080';

  const fetchAssets = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/assets`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setAssets(data || []);
    } catch (e) {
      setError(e.message);
      console.error('Error fetching assets:', e);
    } finally {
      setLoading(false);
    }
  };

  const createAsset = async (newAsset) => {
    try {
      const response = await fetch(`${API_URL}/assets`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newAsset),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      fetchAssets(); // Refresh the list
    } catch (e) {
      console.error('Error creating asset:', e);
      // Optionally re-throw or handle state update for UI
    }
  };

  const updateAsset = async (updatedAsset) => {
    try {
      const response = await fetch(`${API_URL}/assets/${updatedAsset.DEALERID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedAsset),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      fetchAssets(); // Refresh the list
    } catch (e) {
      console.error('Error updating asset:', e);
    }
  };

  const deleteAsset = async (id) => {
    try {
      const response = await fetch(`${API_URL}/assets/${id}`, { method: 'DELETE' });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      fetchAssets(); // Refresh the list
    } catch (e) {
      console.error('Error deleting asset:', e);
    }
  };

  useEffect(() => {
    fetchAssets();
  }, []);

  const value = {
    assets,
    loading,
    error,
    fetchAssets,
    createAsset,
    updateAsset,
    deleteAsset,
  };

  return <AssetContext.Provider value={value}>{children}</AssetContext.Provider>;
};
