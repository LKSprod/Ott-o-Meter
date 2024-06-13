<script lang="ts">
  import type { GrowUnit } from "./types";

  export let growunit: GrowUnit;
</script>

<div class="editor">
  <input bind:value={growunit.Name} />
</div>

<style>
  .editor {
    background: red;
    padding: 5rem;
    border-radius: 2.5rem;
  }
</style>
