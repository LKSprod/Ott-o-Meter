<script lang="ts">
  import GrowunitEditor from "$lib/growunits/GrowunitEditor.svelte";
  import type { GrowUnit } from "$lib/growunits/types";

  let growunit: GrowUnit = {
    Name: "New Grow Unit",
    Id: 0,
    Width: 0,
    Height: 0,
    Depth: 0,
    CarbonFilter: false,
    ActiveIntake: false,
    OuttakeFanThroughputInM3H: 0,
    WattageLamp: 0,
    Ventilation: false,
    Inside: false,
  };
</script>

<GrowunitEditor bind:growunit />
