nclude <Adafruit_NeoPixel.h>

#define LED_PIN 17
#define SENSOR_PIN A0

#define PIXELS 24
#define BRIGHTNESS 50

#define Strip Adafruit_NeoPixel
#define rgb(r, g, b) strip.Color(r, g, b)

enum Command {cmd_on, cmd_off, cmd_animate,
							cmd_sensor, cmd_raw_sensor,
							cmd_set_brightness,
							cmd_set_threshold};

Strip strip(PIXELS, LED_PIN);

int threshold = 0;
uint32_t init_color = rgb(237, 17, 230);
uint32_t off_color = rgb(0, 0, 0);
uint32_t green = rgb(0, 255, 0);
uint32_t red = rgb(255, 0, 0);

void setup() {
	Serial.begin(9600);
	strip.begin();
	strip.setBrightness(BRIGHTNESS);
	turn_off();
	circle_animate(init_color, 20);
	delay(10);
	circle_animate(off_color, 20);
	while (!Serial) delay(10);
	Serial.write(96);
	flash(init_color, 100);
}

void loop() {
	while (Serial.available() == 0) delay(5);
	if (Serial.available() > 0) {
		Command mode = static_cast<Command>(Serial.read());

		int reading;
		int input;
		uint32_t color;
		unsigned long ms;

		switch (mode) {
		case cmd_on:
			color = read_color();
			set_color(color);
			break;
		case cmd_off:
			turn_off();
			break;
		case cmd_animate:
			color = read_color();
			ms = Serial.read();
			circle_animate(color, ms);
			break;
		case cmd_sensor:
			reading = read_sensor();
			Serial.write(reading);
			break;
		case cmd_raw_sensor:
			reading = analogRead(SENSOR_PIN);
			Serial.write(reading);
			break;
		case cmd_set_brightness:
			input = Serial.read();
			set_brightness(input);
			break;
		case cmd_set_threshold:
			input = Serial.read();
			set_threshold(input);
			break;
		}
	}
}

void set_threshold(int value) {
	if (threshold == value) return;
	threshold = value;
	// EEPROM.put(THRESHOLD, value);
}


int read_threshold() {
	return threshold;
}

uint32_t read_color() {
	char buf[3];
	Serial.readBytes(buf, 3);
	int r = buf[0];
	int g = buf[1];
	int b = buf[2];
	return rgb(r, g, b);
}

void set_brightness(uint8_t value) {
	strip.setBrightness(value);
}

void turn_off() {
	strip.clear();
	strip.show();
}
void set_color(uint32_t color) {
	strip.fill(color);
	strip.show();
}

void flash(uint32_t color, unsigned long ms) {
	set_color(color);
	delay(ms);
	turn_off();
}

void circle_animate(uint32_t color, unsigned long ms) {
	int pixels = strip.numPixels();
	for (int i = 0; i < pixels; i++) {
		strip.setPixelColor(i, color);
		strip.show();
		delay(ms);
	}
}

int read_sensor() {
	int value = analogRead(SENSOR_PIN);
	return value < threshold;
}
