<script setup></script>

<template>
  <div>
    <div class="controls">
      <h2>Code Playground</h2>
      <div class="btn-group">
        <button @click="runCode">Run</button>
        <button @click="check">Check</button>
      </div>
    </div>
    <div class="io">
      <textarea
        id="code"
        v-model="inp"
        @keydown="preventTab"
        autocomplete="off"
        name="code"
        spellcheck="false"
        wrap="off"
        ></textarea
      >
      <div class="output">
        <div class="output-header">Program output:</div>
        <pre class="code-output">{{ out }}</pre>
      </div>
    </div>
  </div>
</template>

<script>
const placeholder = `package main

import "fmt"

func main(){
    fmt.Println("hello world")
}`;

export default {
  data() {
    return {
      inp: placeholder,
      out: "",
      running: false,
    };
  },
  methods: {
    preventTab(event) {
      // if not tab, don't care
      if (event.key !== "Tab") {
        return;
      }
      event.preventDefault();
    },
    async check() {
      const res = await fetch("http://localhost:8008/check");
      const r = await res.json();
      console.log(r);
    },
    async runCode() {
      let data = {
        inp: this.inp,
        wheat: true,
      };
      console.log(JSON.stringify(data));

      this.running = true
      const res = await fetch("http://localhost:8008/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const r = await res.json();
      console.log(r);
      this.out = r.Value
      this.running = false
    },
  },
};
</script>

<style>
@import "./assets/base.css";

#app {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;

  font-weight: normal;
}

.controls {
  display: flex;
  flex-direction: row;
  justify-content: space-between;

  max-width: 600px;
  margin: 10px auto;
}

.btn-group {
  display: flex;
  flex-direction: row;
  gap: 2px;
}

button {
  border: none;
  padding: 0.75em 1.5em;
  border-radius: 2px;
}

button:hover {
  cursor: pointer;
  background-color: var(--vt-c-white-soft);
}

.io {
  height: 40em;

  display: flex;
  flex-direction: row;
}

.io textarea {
  height: 100%;
  flex-basis: 70%;
  resize: none;
}

.output{
  background-color: var(--color-background-soft);
  flex-basis: 30%;
  flex-grow: 0;
}
.output-header{
  padding: 1em;
  border-bottom: 1px solid white;
}

.code-output{
  padding: 1em;
  word-break: break-all;
}

</style>
