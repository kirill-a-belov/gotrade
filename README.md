# 🌀 GoTrade: Real-Time Modular Trading Framework for Crypto Futures

**GoTrade** is a high-performance, modular trading framework written in Go (backend) and JavaScript (frontend), designed for **real-time futures trading on Bybit** and easily extensible to any exchange or asset.  
It provides infrastructure for **strategy development, live market monitoring, and automated execution**.

---

## 🚀 Features

### ✅ Real-Time Engine
- Tick-level market data listener (tested with [Bybit](https://www.bybit.com/en/) WebSocket)
- Real-time candle aggregation across **unlimited timeframes**
- Candle-based **live signal calculation** on every tick
- Supports **multiple concurrent signals** per timeframe

### 🧠 Strategy Layer
- Plug in multiple strategies and signals
- Each tick triggers candle update + signal evaluation
- Full market picture in real time from multiple perspectives

### 🛠️ Position Manager
- Smart position control (open/close/adjust based on signals)
- Integrated **risk management**
- Dynamic **trailing stop logic**

### 📊 Web UI
- Built-in web server with live chart rendering
- Displays:
    - Candles and trades
    - Signals on candles
    - Position markers (entry/exit/stop)
    - Signal metadata and trigger details

### 🔮 Roadmap (In Progress)
- Local exchange simulator for safe strategy development
- Historical backtesting with candle reconstruction
- Automatic parameter tuning for strategy optimization

---

## 🧩 Stack
- **Backend:** Go (high performance, concurrency-ready)
- **Frontend:** JavaScript + [TradingView](https://www.tradingview.com/) charts
- **API:** WebSocket for market data + REST for execution

---

## 📷 Screenshot

![GoTrade Preview](./screenshot.png)

---

## 📦 Installation

```bash
git clone https://github.com/kirill-a-belov/gotrade
cd gotrade
go build
./gotrade
```

## Status

Open-source and research-oriented.

Designed for crypto trading strategy experimentation, real-time market analysis, and trading infrastructure research. 

Contributions and forks are welcome.

## Tags

#Go #cryptocurrency #trading #futures #real_time #modular_framework
