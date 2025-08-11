import { useState } from 'react';
import { useAssets } from './context/AssetContext';

function App() {
  const { assets, loading, error, createAsset, updateAsset, deleteAsset } = useAssets();
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

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewAsset({ ...newAsset, [name]: value });
  };

  const handleEditingInputChange = (e) => {
    const { name, value } = e.target;
    setEditingAsset({ ...editingAsset, [name]: value });
  };

  const handleCreateSubmit = (e) => {
    e.preventDefault();
    createAsset(newAsset);
    setNewAsset({
      DEALERID: '', MSISDN: '', MPIN: '', BALANCE: '', STATUS: 'active',
      TRANSAMOUNT: '0', TRANSTYPE: 'init', REMARKS: '',
    });
  };

  const handleUpdateSubmit = (e) => {
    e.preventDefault();
    updateAsset(editingAsset);
    setEditingAsset(null);
  };

  const startEdit = (asset) => {
    setEditingAsset({ ...asset });
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Asset Tracker</h1>

      {error && <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">{error}</div>}

      <div className="mb-8 p-4 border rounded-lg">
        <h2 className="text-xl font-semibold mb-2">{editingAsset ? 'Edit Asset' : 'Create New Asset'}</h2>
        <form onSubmit={editingAsset ? handleUpdateSubmit : handleCreateSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input
            type="text"
            name="DEALERID"
            placeholder="Dealer ID"
            value={editingAsset ? editingAsset.DEALERID : newAsset.DEALERID}
            onChange={editingAsset ? handleEditingInputChange : handleInputChange}
            className="border p-2 rounded"
            disabled={!!editingAsset}
            required
          />
          <input
            type="text"
            name="MSISDN"
            placeholder="MSISDN"
            value={editingAsset ? editingAsset.MSISDN : newAsset.MSISDN}
            onChange={editingAsset ? handleEditingInputChange : handleInputChange}
            className="border p-2 rounded"
            required
          />
          <input
            type="text"
            name="MPIN"
            placeholder="MPIN"
            value={editingAsset ? editingAsset.MPIN : newAsset.MPIN}
            onChange={editingAsset ? handleEditingInputChange : handleInputChange}
            className="border p-2 rounded"
            required
          />
          <input
            type="text"
            name="BALANCE"
            placeholder="Balance"
            value={editingAsset ? editingAsset.BALANCE : newAsset.BALANCE}
            onChange={editingAsset ? handleEditingInputChange : handleInputChange}
            className="border p-2 rounded"
            required
          />
          <input
            type="text"
            name="REMARKS"
            placeholder="Remarks"
            value={editingAsset ? editingAsset.REMARKS : newAsset.REMARKS}
            onChange={editingAsset ? handleEditingInputChange : handleInputChange}
            className="border p-2 rounded md:col-span-2"
          />
          <div className="md:col-span-2 flex justify-end gap-2">
            {editingAsset && (
              <button
                type="button"
                onClick={() => setEditingAsset(null)}
                className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded"
              >
                Cancel Edit
              </button>
            )}
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
            >
              {editingAsset ? 'Update Asset' : 'Create Asset'}
            </button>
          </div>
        </form>
      </div>

      <div>
        <h2 className="text-xl font-semibold mb-2">Assets on Ledger</h2>
        {loading ? (
          <p>Loading assets...</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead className="bg-gray-200">
                <tr>
                  <th className="p-2 border">Dealer ID</th>
                  <th className="p-2 border">MSISDN</th>
                  <th className="p-2 border">Balance</th>
                  <th className="p-2 border">Status</th>
                  <th className="p-2 border">Remarks</th>
                  <th className="p-2 border">Actions</th>
                </tr>
              </thead>
              <tbody>
                {assets.map((asset) => (
                  <tr key={asset.DEALERID} className="hover:bg-gray-100">
                    <td className="p-2 border">{asset.DEALERID}</td>
                    <td className="p-2 border">{asset.MSISDN}</td>
                    <td className="p-2 border">{asset.BALANCE}</td>
                    <td className="p-2 border">{asset.STATUS}</td>
                    <td className="p-2 border">{asset.REMARKS}</td>
                    <td className="p-2 border text-center">
                      <button
                        onClick={() => startEdit(asset)}
                        className="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-1 px-2 rounded mr-2"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => deleteAsset(asset.DEALERID)}
                        className="bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-2 rounded"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
