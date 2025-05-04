# Automatic Plant Watering


## Using RP2040

## Using DAPLINK instead of STM-LINK v2

Create `bluepill-dap.json` in `tinygo/targets`:

```json
{
	"inherits": ["cortex-m3"],
	"build-tags": ["bluepill", "stm32f103", "stm32f1", "stm32"],
	"serial": "uart",
	"linkerscript": "targets/stm32.ld",
	"extra-files": [
		"src/device/stm32/stm32f103.s"
	],
	"flash-method": "openocd",
	"openocd-interface": "cmsis-dap",
	"openocd-target": "stm32f1x"
}
```

and flash with `tinygo flash --target bluepill-dap`.

NOTE: Bluepill I had did not read ADCs at all - the function call hang.