package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon

#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>
#import <ApplicationServices/ApplicationServices.h>
#include <stdlib.h>

extern char *cliperPopupItems(void);
extern void cliperPopupItemClicked(int index);

@interface CliperPopupTarget : NSObject
- (void)press:(id)sender;
@end

@implementation CliperPopupTarget
- (void)press:(id)sender {
	NSNumber *index = [sender representedObject];
	cliperPopupItemClicked(index.intValue);
}
@end

static EventHotKeyRef cliperHotKeyRef = NULL;
static CliperPopupTarget *cliperPopupTarget = NULL;

static bool cliperRequestEventAccess(void) {
	if (@available(macOS 10.15, *)) {
		if (CGPreflightPostEventAccess()) {
			return true;
		}
		return CGRequestPostEventAccess();
	}
	return true;
}

static void cliperPaste(void) {
	dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 120 * NSEC_PER_MSEC), dispatch_get_main_queue(), ^{
		CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
		if (source == NULL) {
			return;
		}

		CGEventRef keyDown = CGEventCreateKeyboardEvent(source, (CGKeyCode)kVK_ANSI_V, true);
		CGEventRef keyUp = CGEventCreateKeyboardEvent(source, (CGKeyCode)kVK_ANSI_V, false);
		if (keyDown == NULL || keyUp == NULL) {
			if (keyDown != NULL) {
				CFRelease(keyDown);
			}
			if (keyUp != NULL) {
				CFRelease(keyUp);
			}
			CFRelease(source);
			return;
		}

		CGEventSetFlags(keyDown, kCGEventFlagMaskCommand);
		CGEventSetFlags(keyUp, kCGEventFlagMaskCommand);
		CGEventPost(kCGHIDEventTap, keyDown);
		CGEventPost(kCGHIDEventTap, keyUp);

		CFRelease(keyDown);
		CFRelease(keyUp);
		CFRelease(source);
	});
}

static void cliperShowPopupMenu(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		char *json = cliperPopupItems();
		if (json == NULL) {
			return;
		}

		NSString *jsonString = [NSString stringWithUTF8String:json];
		free(json);
		if (jsonString == nil) {
			return;
		}

		NSData *data = [jsonString dataUsingEncoding:NSUTF8StringEncoding];
		NSArray *items = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
		if (![items isKindOfClass:[NSArray class]]) {
			return;
		}

		if (cliperPopupTarget == nil) {
			cliperPopupTarget = [CliperPopupTarget new];
		}

		NSMenu *menu = [NSMenu new];
		menu.autoenablesItems = false;
		for (NSDictionary *dict in items) {
			if ([dict[@"separator"] boolValue]) {
				[menu addItem:[NSMenuItem separatorItem]];
				continue;
			}

			NSString *text = dict[@"text"] ?: @"";
			NSMenuItem *item = [menu addItemWithTitle:text action:nil keyEquivalent:@""];
			BOOL enabled = [dict[@"enabled"] boolValue];
			item.enabled = enabled;
			if (enabled) {
				item.target = cliperPopupTarget;
				item.action = @selector(press:);
				item.representedObject = dict[@"index"];
			}
		}

		NSPoint mouseLocation = [NSEvent mouseLocation];
		[menu popUpMenuPositioningItem:nil atLocation:mouseLocation inView:nil];
	});
}

static OSStatus cliperHotKeyHandler(EventHandlerCallRef nextHandler, EventRef event, void *userData) {
	cliperShowPopupMenu();
	return noErr;
}

static int cliperStartHotkey(void) {
	if (cliperHotKeyRef != NULL) {
		return noErr;
	}

	EventTypeSpec eventType;
	eventType.eventClass = kEventClassKeyboard;
	eventType.eventKind = kEventHotKeyPressed;
	OSStatus handlerStatus = InstallApplicationEventHandler(&cliperHotKeyHandler, 1, &eventType, NULL, NULL);
	if (handlerStatus != noErr) {
		return handlerStatus;
	}

	EventHotKeyID hotKeyID;
	hotKeyID.signature = 'CLPV';
	hotKeyID.id = 1;
	return RegisterEventHotKey(kVK_ANSI_V, optionKey, hotKeyID, GetApplicationEventTarget(), 0, &cliperHotKeyRef);
}
*/
import "C"
import "log"

func startHotkey(app *App) {
	activeApp = app
	if !bool(C.cliperRequestEventAccess()) {
		log.Print("Cliper needs Accessibility permission to paste selected clipboard items")
	}
	if status := C.cliperStartHotkey(); status != 0 {
		log.Printf("failed to register Option+V hotkey: %d", int(status))
	}
}

func pasteFromHotkeyMenu() {
	C.cliperPaste()
}
