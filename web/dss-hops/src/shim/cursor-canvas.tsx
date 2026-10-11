/**
 * Blog web host shim: extends @thisismydesign/cursor-canvas-web for DSS hops canvas.
 */
import {
  Button as MantineButton,
  Code as MantineCode,
  Collapse,
  Table as MantineTable,
  Title,
} from "@mantine/core";
import type { CSSProperties, ReactNode } from "react";
import {
  useEffect,
  useRef,
  useState,
  useCallback,
  useMemo,
} from "react";
import * as base from "@thisismydesign/cursor-canvas-web";
import { useHostTheme as useBaseHostTheme } from "@thisismydesign/cursor-canvas-web";

export { useEffect, useRef, useState };

export const {
  H1,
  H2,
  Text,
  Stack,
  Row,
  Spacer,
  Grid,
  Divider,
  Card,
  CardHeader,
  CardBody,
  Stat,
  Pill,
  Callout,
  Select,
  Checkbox,
  TextInput,
  LineChart,
  useCanvasState,
} = base;

const category = {
  blue: "#228be6",
  orange: "#fd7e14",
  green: "#40c057",
  red: "#fa5252",
  yellow: "#fab005",
  violet: "#7950f2",
  gray: "#868e96",
};

export function useHostTheme() {
  const t = useBaseHostTheme();
  return useMemo(
    () => ({
      ...t,
      category,
      fill: {
        primary: "rgba(20, 20, 20, 0.06)",
        secondary: "rgba(20, 20, 20, 0.04)",
        tertiary: "rgba(20, 20, 20, 0.08)",
        quaternary: "rgba(20, 20, 20, 0.03)",
      },
      stroke: {
        primary: "rgba(20, 20, 20, 0.2)",
        secondary: "rgba(20, 20, 20, 0.14)",
        tertiary: "rgba(20, 20, 20, 0.1)",
      },
      text: {
        ...t.text,
        onAccent: "#ffffff",
      },
    }),
    [t],
  );
}

export function H3({
  children,
  style,
}: {
  children?: ReactNode;
  style?: CSSProperties;
}) {
  return (
    <Title order={3} style={style}>
      {children}
    </Title>
  );
}

export function Button({
  children,
  onClick,
  style,
}: {
  children?: ReactNode;
  onClick?: () => void;
  style?: CSSProperties;
}) {
  return (
    <MantineButton variant="light" size="xs" onClick={onClick} style={style}>
      {children}
    </MantineButton>
  );
}

export function Code({
  children,
  style,
}: {
  children?: ReactNode;
  style?: CSSProperties;
}) {
  return <MantineCode style={style}>{children}</MantineCode>;
}

type TableRowTone = "info" | undefined;

export function Table({
  headers,
  rows,
  columnAlign,
  rowTone,
  style,
}: {
  headers: ReactNode[];
  rows: ReactNode[][];
  columnAlign?: Array<"left" | "right" | "center">;
  rowTone?: Array<TableRowTone>;
  style?: CSSProperties;
}) {
  return (
    <MantineTable striped highlightOnHover withTableBorder style={style}>
      <MantineTable.Thead>
        <MantineTable.Tr>
          {headers.map((h, idx) => (
            <MantineTable.Th key={idx}>{h}</MantineTable.Th>
          ))}
        </MantineTable.Tr>
      </MantineTable.Thead>
      <MantineTable.Tbody>
        {rows.map((row, i) => (
          <MantineTable.Tr
            key={i}
            style={
              rowTone?.[i] === "info"
                ? { backgroundColor: "rgba(34, 139, 230, 0.08)" }
                : undefined
            }
          >
            {row.map((cell, j) => (
              <MantineTable.Td
                key={j}
                style={{
                  textAlign: columnAlign?.[j] ?? "left",
                }}
              >
                {cell}
              </MantineTable.Td>
            ))}
          </MantineTable.Tr>
        ))}
      </MantineTable.Tbody>
    </MantineTable>
  );
}

export type { TableRowTone };

export function CollapsibleSection({
  title,
  children,
  defaultOpen,
  count,
}: {
  title: string;
  children?: ReactNode;
  defaultOpen?: boolean;
  count?: number;
}) {
  const [open, setOpen] = useState(defaultOpen ?? false);
  const label = count !== undefined ? `${title} (${count})` : title;
  return (
    <div style={{ marginTop: 8 }}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        style={{
          background: "none",
          border: "none",
          cursor: "pointer",
          fontWeight: 600,
          padding: 0,
          textAlign: "left",
        }}
      >
        {open ? "▼ " : "▶ "}
        {label}
      </button>
      <Collapse in={open}>{children}</Collapse>
    </div>
  );
}

// Satisfy re-exports the canvas might expect from the SDK surface
export { useCallback, useMemo };
