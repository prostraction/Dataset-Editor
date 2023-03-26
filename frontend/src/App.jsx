import { useState } from "react";
import "./App.css";
import { SetDirectoryDialog } from "../wailsjs/go/main/App";
import { StartMergeProcess } from "../wailsjs/go/main/App";

function App() {
  //const [resultText, setResultText] = useState(
  //  "Please enter your name below 👇"
  //);
  //const [name, setName] = useState("");
  //const updateName = (e) => setName(e.target.value);
  //const updateResultText = (result) => setResultText(result);

  const [dir1, dir1set] = useState("Select DIR1");
  const updateDir1Text = (selectdir1) => dir1set(selectdir1);
  const [dir2, dir2set] = useState("Select DIR2");
  const updateDir2Text = (selectdir2) => dir2set(selectdir2);

  //function greet() {
  //  Greet(name).then(updateResultText);
  //}

  function SetDirectory1() {
    SetDirectoryDialog().then(updateDir1Text);
  }
  function SetDirectory2() {
    SetDirectoryDialog().then(updateDir2Text);
  }
  function CallMergeProgress() {
    StartMergeProcess("test", "test", "test")
  }
  //<input id="dirSelected" directory="" webkitdirectory="" type="file" />
  return (
    <div id="App">
      <div class="selectDirField">
        <text id="selectdir1" className="selectdir1">
          {dir1}
        </text>
        <button className="btn" onClick={SetDirectory1}>
          Select First
        </button>
      </div>

      <div class="selectDirField">
        <text id="selectdir2" className="selectdir2">
          {dir2}
        </text>
        <button className="btn" onClick={SetDirectory2}>
          Select Second
        </button>
      </div>
    </div>
  );
}

export default App;
