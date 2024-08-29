[繁體中文](https://github.com/TimLai666/Enhanced-Python-Compiler?tab=readme-ov-file#%E5%BC%B7%E5%8C%96%E7%89%88-python-%E7%B7%A8%E8%AD%AF%E5%99%A8%E5%B0%88%E6%A1%88)
### Enhanced Python Compiler Project

---

#### **Project Overview**

The Enhanced Python Compiler Project aims to create a hybrid compiler that combines the strengths of Go's performance with Python's flexibility. The project is divided into two primary components: an enhanced runtime implemented in Go and the packaging logic that bundles the custom runtime, CPython interpreter, and Python code into a single executable.

---

### **Architecture**

#### **1. Enhanced Runtime (Implemented in Go)**

- **Purpose:**  
  The enhanced runtime is responsible for optimizing Python code execution by converting Python's dynamic constructs into more efficient, Go-native structures. This includes:
  - **Type Inference:** Analyzing variables at runtime to determine their types, converting them into Go's static types when possible.
  - **Constant Folding:** Detecting variables that don't change throughout their lifetime and converting them into constants to enhance performance.
  - **Efficient Data Structures:** Replacing Python's standard data structures with more performant Go equivalents when beneficial.

- **Execution Strategy:**  
  The first time Python code is run, the runtime will perform necessary conversions and optimizations. The result will be cached or saved as a binary, eliminating the need for repeated dynamic interpretation in subsequent executions.

- **JIT-like and AOT-like Behavior:**
  - **First Execution:** Behaves like Just-In-Time (JIT) compilation, where Python code is analyzed and optimized during the first run.
  - **Subsequent Executions:** The code is executed as an Ahead-Of-Time (AOT) compiled binary, ensuring optimal performance without the overhead of dynamic interpretation.

#### **2. Packaging Logic**

- **Purpose:**  
  The packaging logic bundles the optimized Python code, the enhanced runtime, and the CPython interpreter into a single, standalone executable.

- **Bundling Process:**  
  - **Python Code:** The original Python source or its optimized form is included.
  - **Enhanced Runtime:** The Go-based runtime responsible for executing and optimizing the Python code.
  - **CPython Interpreter:** Provides compatibility with Python's standard libraries and ensures the execution of Python code that relies on the CPython ecosystem.
  - **Dependencies:** All required Python packages and dependencies are bundled to ensure that the executable is self-contained and does not require any external Python installation.

---

### **Project Structure**

```
enhanced_python_compiler/
│
├── cmd/
│   └── compiler/
│       ├── main.go               # Entry point for the compiler
│       └── ...                   # Other Go files related to CLI operations
│
├── internal/
│   ├── runtime/
│   │   ├── runtime.go            # Go-based enhanced runtime implementation
│   │   ├── type_inference.go     # Logic for type inference
│   │   ├── constant_folding.go   # Logic for constant folding
│   │   └── ...                   # Additional runtime optimization files
│   │
│   ├── translator/
│   │   ├── translator.go         # Translates Python AST into optimized Go structures
│   │   └── ...                   # Additional translation and optimization files
│   │
│   ├── packager/
│   │   ├── packager.go           # Logic for bundling runtime, CPython, and Python code
│   │   └── dependency_resolver.go# Resolves and bundles Python dependencies
│   │
│   └── parser/
│       ├── parser.go             # Parses Python code into AST
│       └── ...                   # Additional parsing and AST manipulation files
│
└── README.md                     # This file
```

---

### **How It Works**

1. **Code Parsing:**  
   The Python code is first parsed into an Abstract Syntax Tree (AST) using the `parser` package.

2. **Translation and Optimization:**  
   The AST is passed to the `translator`, where it is translated into Go-native structures. During this process, the `runtime` package applies various optimizations like type inference and constant folding.

3. **Packaging:**  
   The optimized code, along with the custom runtime and CPython interpreter, is bundled into a single executable using the `packager`.

4. **Execution:**  
   The resulting executable can be run directly, where it will leverage the optimizations made during the first execution, providing AOT-like performance.

---

### **Future Work**

- **Further Optimization:**  
  Continue refining the runtime to support more complex Python constructs and optimize them effectively.
  
- **Enhanced Compatibility:**  
  Improve the bundling process to ensure the widest compatibility with various Python packages and versions.

- **Community and Documentation:**  
  Expand this README and other documentation to facilitate contributions and clarify usage for other developers.

---

---

### 強化版 Python 編譯器專案

---

#### **專案概述**

強化版 Python 編譯器專案旨在創建一個混合型編譯器，結合 Go 的效能優勢與 Python 的靈活性。專案分為兩個主要部分：由 Go 實現的增強版執行時和將自訂執行時、CPython 解釋器及 Python 程式碼打包成單一可執行檔的打包邏輯。

---

### **架構設計**

#### **1. 增強版執行時（由 Go 實現）**

- **目的：**  
  增強版執行時負責通過將 Python 的動態結構轉換為更高效的 Go 原生結構來優化 Python 程式碼的執行，包括：
  - **型別推斷：** 分析執行時的變數以確定其型別，並在可能的情況下將其轉換為 Go 的靜態型別。
  - **常數折疊：** 檢測在其生命週期內未改變的變數，並將其轉換為常數以提高效能。
  - **高效資料結構：** 當有利於性能時，將 Python 標準資料結構替換為性能更佳的 Go 等效結構。

- **執行策略：**  
  在第一次執行 Python 程式碼時，執行時將進行必要的轉換和優化。結果將被快取或保存為二進制檔，從而消除在後續執行中反覆進行動態解釋的需要。

- **JIT 與 AOT 行為：**
  - **首次執行：** 表現得像 Just-In-Time (JIT) 編譯，在首次執行期間分析和優化 Python 程式碼。
  - **後續執行：** 程式碼以 Ahead-Of-Time (AOT) 編譯的二進制形式執行，確保最佳效能，無需動態解釋的開銷。

#### **2. 打包邏輯**

- **目的：**  
  打包邏輯將優化過的 Python 程式碼、增強版執行時和 CPython 解釋器打包成一個獨立的可執行檔。

- **打包過程：**  
  - **Python 程式碼：** 包含原始的 Python 原始碼或其優化形式。
  - **增強版執行時：** 由 Go 實現的執行時，負責執行和優化 Python 程式碼。
  - **CPython 解釋器：** 提供與 Python 標準庫的相容性，確保依賴於 CPython 生態系的 Python 程式碼能夠正常執行。
  - **依賴套件：** 所有所需的 Python 套件和依賴項將被打包，以確保可執行檔是自包含的，無需外部的 Python 安裝。

---

### **專案結構**

```
enhanced_python_compiler/
│
├── cmd/
│   └── compiler/
│       ├── main.go               # 編譯器的進入點
│       └── ...                   # 其他與 CLI 操作相關的 Go 檔案
│
├── internal/
│   ├── runtime/
│   │   ├── runtime.go            # 由 Go 實現的增強版執行時
│   │   ├── type_inference.go     # 型別推斷邏輯
│   │   ├── constant_folding.go   # 常數折疊邏輯
│   │   └── ...                   # 其他執行時優化相關的檔案
│   │
│   ├── translator/
│   │   ├── translator.go         # 將 Python AST 轉換為優化的 Go 結構
│   │   └── ...                   # 其他翻譯和優化相關的檔案
│   │
│   ├── packager/
│   │   ├── packager.go           # 負責打包執行時、CPython 和 Python 程式碼
│   │   └── dependency_resolver.go# 解決並打包 Python 依賴項
│   │
│   └── parser

/
│       ├── parser.go             # 將 Python 程式碼解析為 AST
│       └── ...                   # 其他解析和 AST 操作相關的檔案
│
└── README.md                     # 此檔案
```

---

### **運作方式**

1. **程式碼解析：**  
   使用 `parser` 套件將 Python 程式碼首先解析為抽象語法樹（AST）。

2. **翻譯與優化：**  
   AST 被傳遞給 `translator`，該模組將其轉換為 Go 原生結構。在此過程中，`runtime` 套件應用各種優化，如型別推斷和常數折疊。

3. **打包：**  
   優化過的程式碼，連同自訂的執行時和 CPython 解釋器，被 `packager` 打包成一個可執行檔。

4. **執行：**  
   最終的可執行檔可以直接運行，並利用首次執行時所做的優化，提供 AOT 類似的性能。

---

### **未來工作**

- **進一步優化：**  
  繼續改進執行時，支持更複雜的 Python 結構並有效地優化它們。
  
- **增強相容性：**  
  改進打包過程，以確保與各種 Python 套件和版本的廣泛相容性。

- **社群與文件：**  
  擴展此 README 和其他文檔，以促進貢獻並澄清其他開發者的使用情境。
