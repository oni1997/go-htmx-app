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
        } catch (error) {
            console.error('Error connecting to MetaMask', error);
            updateConnectionStatus(false);
        }
    } else {
        alert('MetaMask is not installed. Please install it to use this app.');
        updateConnectionStatus(false);
    }
}

function updateConnectionStatus(isConnected) {
    const statusElement = document.getElementById('connection-status');
    if (isConnected) {
        statusElement.textContent = '🟢 Connected';
        statusElement.classList.add('connected');
    } else {
        statusElement.textContent = '⚪ Not connected';
        statusElement.classList.remove('connected');
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
            document.getElementById('balance').textContent = `Balance: ${parseFloat(balanceInCUSD).toFixed(6)} cUSD`;
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

// Make sure this event listener is set up
document.getElementById('close-modal').addEventListener('click', closeModal);
document.getElementById('transfer-form').addEventListener('submit', transferCUSD);

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
    ethereum.on('accountsChanged', connectWallet);
}

// Initial connection attempt
connectWallet();