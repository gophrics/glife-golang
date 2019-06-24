import ProfileService
import TravelService
import Context

def main():
    ProfileService.username_exist()
    success = ProfileService.register()
    if success:
        ProfileService.get_me()
        ProfileService.get_user()
        TravelService.save_trip()
        TravelService.get_all_trips(Context.username)

if __name__ == "__main__":
    main()
