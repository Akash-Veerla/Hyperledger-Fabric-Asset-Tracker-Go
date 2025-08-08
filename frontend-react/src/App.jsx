import { useState, useEffect } from 'react';

function App() {
  const [assets, setAssets] = useState([]);
  const [newAsset, setNewAsset] = useState({
    DEALERID: '',
    MSISDN: '',
    MPIN: '',
    BALANCE: '',
    STATUS: 'active',
    TRANSAMOUNT: '0',
    TRANSTYPE: 'init',
    REMARKS: '',
  });
  const [editingAsset, setEditingAsset] = useState(null);

  const API_URL = 'http://localhost:8080';

  useEffect(() => {
    fetchAssets();
  }, []);

  const fetchAssets = async () => {
    try {
      const response = await fetch(`${API_URL}/assets`);
      const data = await response.json();
      setAssets(data || []);
    } catch (error) {
      console.error('Error fetching assets:', error);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewAsset({ ...newAsset, [name]: value });
  };

  const handleEditingInputChange = (e) => {
    const { name, value } = e.target;
    setEditingAsset({ ...editingAsset, [name]: value });
  };

  const createAsset = async (e) => {
    e.preventDefault();
    try {
      await fetch(`${API_URL}/assets`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newAsset),
      });
      fetchAssets();
      setNewAsset({
        DEALERID: '',
        MSISDN: '',
        MPIN: '',
        BALANCE: '',
        STATUS: 'active',
        TRANSAMOUNT: '0',
        TRANSTYPE: 'init',
        REMARKS: '',
      });
    } catch (error) {
      console.error('Error creating asset:', error);
    }
  };

  const updateAsset = async (e) => {
    e.preventDefault();
    try {
      await fetch(`${API_URL}/assets/${editingAsset.DEALERID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editingAsset),
      });
      fetchAssets();
      setEditingAsset(null);
    } catch (error) {
      console.error('Error updating asset:', error);
    }
  };

  const deleteAsset = async (id) => {
    try {
      await fetch(`${API_URL}/assets/${id}`, { method: 'DELETE' });
      fetchAssets();
    } catch (error) {
      console.error('Error deleting asset:', error);
    }
  };


  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Asset Tracker</h1>

      <div className="mb-8">
        <h2 className="text-xl font-semibold mb-2">{editingAsset ? 'Edit Asset' : 'Create New Asset'}</h2>
        <form onSubmit={editingAsset ? updateAsset : createAsset} className="grid grid-cols-2 gap-4">
          <input
            type="text"
            name="DEALERID"
            placeholder="Dealer ID"
            value={editingAsset ? editingAsset.DEALERID : newAsset.DEALERID}
            onChange={editingAsset ? handleEditingInputChange: handleInputChange}
            className="border p-2"
            disabled={editingAsset}
          />
          <input
            type="text"
            name="MSISDN"
            placeholder="MSISDN"
            value={editingAsset ? editingAsset.MSISDN : newAsset.MSISDN}
            onChange={editingAsset ? handleEditingInputChange: handleInputChange}
            className="border p-2"
          />
          <input
            type="text"
            name="MPIN"
            placeholder="MPIN"
            value={editingAsset ? editingAsset.MPIN : newAsset.MPIN}
            onChange={editingAsset ? handleEditingInputChange: handleInputChange}
            className="border p-2"
          />
          <input
            type="text"
            name="BALANCE"
            placeholder="Balance"
            value={editingAsset ? editingAsset.BALANCE : newAsset.BALANCE}
            onChange={editingAsset ? handleEditingInputChange: handleInputChange}
            className="border p-2"
          />
          <input
            type="text"
            name="REMARKS"
            placeholder="Remarks"
            value={editingAsset ? editingAsset.REMARKS : newAsset.REMARKS}
            onChange={editingAsset ? handleEditingInputChange: handleInputChange}
            className="border p-2 col-span-2"
          />
          <button type="submit" className="bg-blue-500 text-white p-2 rounded col-span-2">
            {editingAsset ? 'Update Asset' : 'Create Asset'}
          </button>
          {editingAsset && (
            <button
              type="button"
              onClick={() => setEditingAsset(null)}
              className="bg-gray-500 text-white p-2 rounded col-span-2"
            >
              Cancel Edit
            </button>
          )}
        </form>
      </div>

      <div>
        <h2 className="text-xl font-semibold mb-2">Assets</h2>
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-200">
              <th className="border p-2">Dealer ID</th>
              <th className="border p-2">MSISDN</th>
              <th className="border p-2">Balance</th>
              <th className="border p-2">Status</th>
              <th className="border p-2">Remarks</th>
              <th className="border p-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {assets.map((asset) => (
              <tr key={asset.DEALERID}>
                <td className="border p-2">{asset.DEALERID}</td>
                <td className="border p-2">{asset.MSISDN}</td>
                <td className="border p-2">{asset.BALANCE}</td>
                <td className="border p-2">{asset.STATUS}</td>
                <td className="border p-2">{asset.REMARKS}</td>
                <td className="border p-2">
                  <button
                    onClick={() => setEditingAsset(asset)}
                    className="bg-yellow-500 text-white p-1 rounded mr-2"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => deleteAsset(asset.DEALERID)}
                    className="bg-red-500 text-white p-1 rounded"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
