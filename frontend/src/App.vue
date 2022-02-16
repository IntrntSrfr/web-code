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
      if(this.running){return}
      this.running = true
      this.out = 'running..'

      let weed = setInterval(() => {
        this.out += '.'
      }, 500);

      let data = {
        inp: this.inp,
        wheat: true,
      };

      const res = await fetch("http://localhost:8008/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const r = await res.json();

      clearInterval(weed)
      
      this.out = r.Value
      this.running = false
    },
  },
};
</script>

<style>
@import "./assets/base.css";

#app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;

  font-weight: normal;
}

.controls {
  display: flex;
  flex-direction: row;
  justify-content: space-between;

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
  background-color: var(--color-background-black);
  color: var(--color-heading);
  flex-grow: 1;
  resize: none;
}

.io textarea:focus-visible{
  outline: none;
}

.output{
  flex: 0 0 20em;
  background-color: var(--color-background-soft);
  flex-grow: 0;
  overflow-x: hidden;

  display: flex;
  flex-direction: column;
}

.output-header{
  font-size: 1.25em;
  padding: 1em;
  border-bottom: 1px solid var(--color-border);
}

.code-output{
  padding: 1em;
  word-break: break-all;
  overflow-x: scroll;
  flex-grow: 1;
}

</style>
