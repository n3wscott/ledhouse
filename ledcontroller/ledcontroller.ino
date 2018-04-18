#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define PIN 6

Adafruit_NeoPixel strip = Adafruit_NeoPixel(10, PIN, NEO_GRB + NEO_KHZ800);

String inputString = "";
boolean stringComplete = false;

void setup() {
  strip.begin();
  strip.show(); // Initialize all pixels to 'off'

  Serial.begin(115200);
  inputString.reserve(200);
}

void loop() {
  if (stringComplete) {
    Serial.println(inputString);

    if (inputString.startsWith("LED") && inputString.length() == 12) {
      int index = StrToHex(inputString.substring(3,5));
      int red = StrToHex(inputString.substring(5,7));
      int green = StrToHex(inputString.substring(7,9));
      int blue = StrToHex(inputString.substring(9,11));

      Serial.println(index);
      Serial.println(red);
      Serial.println(green);
      Serial.println(blue);

      strip.setPixelColor(index, strip.Color(red, blue, green));
      strip.show();
    } else if (inputString.startsWith("CLEAR")) {
      colorWipe(strip.Color(0, 0, 0), 50);
    }
    
    // clear the string:
    inputString = "";
    stringComplete = false;
  }
}

int StrToHex(String str) {
  return (int) strtol(str.c_str(), 0, 16);
}

void serialEvent() {
  while (Serial.available()) {
    char inChar = (char)Serial.read();
    inputString += inChar;
    if (inChar == '\n') {
      stringComplete = true;
    }
  }
}

void colorWipe(uint32_t c, uint8_t wait) {
  for(uint16_t i=0; i < strip.numPixels(); i++) {
    strip.setPixelColor(i, c);
    strip.show();
    delay(wait);
  }
}
