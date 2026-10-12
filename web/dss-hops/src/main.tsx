import { mountCanvas } from "@thisismydesign/cursor-canvas-web/runtime";
import "@mantine/core/styles.css";
import "@mantine/charts/styles.css";
import DSSCoreTopology from "./DSS-core-topology.canvas";

mountCanvas("root", <DSSCoreTopology />, { defaultColorScheme: "light" });
