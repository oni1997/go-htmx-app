let web3;
let currentAccount;

const cUSDAddress = '0x874069Fa1Eb16D44d622F2e0Ca25eeA172369bC1';

async function connectWallet() {
    if (typeof window.ethereum !== 'undefined') {
        try {
            await ethereum.request({ method: 'eth_requestAccounts' });
            web3 = new Web3(window.ethereum);
            const accounts = await web3.eth.getAccounts();
            currentAccount = accounts[0];
            updateConnectionStatus(true);
            updateBalance();
            document.getElementById('connect-button').style.display = 'none';
            document.getElementById('disconnect-button').style.display = 'block';
        } catch (error) {
            console.error('Error connecting to MetaMask', error);
            updateConnectionStatus(false);
        }
    } else {
        alert('MetaMask is not installed. Please install it to use this app.');
        updateConnectionStatus(false);
    }
}

function disconnectWallet() {
    currentAccount = null;
    updateConnectionStatus(false);
    document.getElementById('connect-button').style.display = 'block';
    document.getElementById('disconnect-button').style.display = 'none';
}

function updateConnectionStatus(isConnected) {
    const statusElement = document.getElementById('connection-status');
    if (isConnected) {
        statusElement.textContent = 'Wallet connected';
    } else {
        statusElement.textContent = 'Wallet disconnected';
    }
}

async function updateBalance() {
    if (currentAccount) {
        try {
            const contract = new web3.eth.Contract([
                {
                    "constant": true,
                    "inputs": [{"name": "_owner", "type": "address"}],
                    "name": "balanceOf",
                    "outputs": [{"name": "balance", "type": "uint256"}],
                    "type": "function"
                }
            ], cUSDAddress);
            
            const balance = await contract.methods.balanceOf(currentAccount).call();
            const balanceInCUSD = web3.utils.fromWei(balance, 'ether');
            document.getElementById('transfer-amount').placeholder = `Amount to Transfer (Max: ${parseFloat(balanceInCUSD).toFixed(6)} cUSD)`;
        } catch (error) {
            console.error('Error fetching balance:', error);
        }
    }
}

async function transferCUSD(event) {
    event.preventDefault();
    const receiverAddress = document.getElementById('receiver-address').value;
    const amount = document.getElementById('transfer-amount').value;

    if (!web3.utils.isAddress(receiverAddress)) {
        alert('Invalid receiver address');
        return;
    }

    const amountInWei = web3.utils.toWei(amount, 'ether');

    try {
        const cUSDContract = new web3.eth.Contract([
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_to",
                        "type": "address"
                    },
                    {
                        "name": "_value",
                        "type": "uint256"
                    }
                ],
                "name": "transfer",
                "outputs": [
                    {
                        "name": "",
                        "type": "bool"
                    }
                ],
                "type": "function"
            }
        ], cUSDAddress);

        const result = await cUSDContract.methods.transfer(receiverAddress, amountInWei).send({
            from: currentAccount
        });

        showSuccessModal(result.transactionHash);
        updateBalance();
    } catch (error) {
        console.error('Error transferring cUSD', error);
        alert('Transaction failed: ' + error.message);
    }
}

function showSuccessModal(transactionHash) {
    const modal = document.getElementById('success-modal');
    const hashElement = document.getElementById('transaction-hash');
    hashElement.textContent = `Transaction Hash: ${transactionHash}`;
    modal.classList.add('show');
}

function closeModal() {
    const modal = document.getElementById('success-modal');
    modal.classList.remove('show');
}

function updateCurrentTime() {
    const now = new Date();
    const timeString = now.toLocaleTimeString('en-US', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' });
    document.getElementById('current-time').textContent = timeString;
}

// Event listeners
document.getElementById('close-modal').addEventListener('click', closeModal);
document.getElementById('transfer-form').addEventListener('submit', transferCUSD);
document.getElementById('connect-button').addEventListener('click', connectWallet);
document.getElementById('disconnect-button').addEventListener('click', disconnectWallet);

// Update time every second
setInterval(updateCurrentTime, 1000);

// Initial time update
updateCurrentTime();

// Check if already connected
if (typeof window.ethereum !== 'undefined') {
    ethereum.request({ method: 'eth_accounts' }).then(accounts => {
        if (accounts.length > 0) {
            connectWallet();
        }
    });
}

// Listen for account changes
if (typeof window.ethereum !== 'undefined') {
    ethereum.on('accountsChanged', (accounts) => {
        if (accounts.length > 0) {
            connectWallet();
        } else {
            disconnectWallet();
        }
    });
}